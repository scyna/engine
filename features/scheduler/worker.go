package scheduler

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
)

type worker struct {
}

func NewWorker() *worker {
	return &worker{
		// qCheck: qb.Insert(scyna_const.DOING_TABLE).
		// 	Columns("bucket", "task_id").
		// 	Unique().
		// 	TTL(60 * time.Second).
		// 	Query(scyna.DB),
		// qGet: qb.Select(scyna_const.TASK_TABLE).
		// 	Columns("id", "topic", "data", "next", "interval", "loop_index", "loop_count", "done").
		// 	Where(qb.Eq("id")).
		// 	Limit(1).
		// 	Query(scyna.DB),
		// qTodos: qb.Select(scyna_const.TODO_TABLE).
		// 	Columns("task_id").
		// 	Where(qb.Eq("bucket")).
		// 	Limit(20).
		// 	Query(scyna.DB),
	}
}

func (w *worker) Start(delay time.Duration, interval time.Duration) {
	go func() {
		time.Sleep(delay)
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				w.execute()
			}
		}
	}()
}

func (w *worker) execute() {
	bucket := GetBucket(time.Now())
	for {
		var tasks []int64
		// qTodos: qb.Select(scyna_const.TODO_TABLE).
		// Columns("task_id").
		// Where(qb.Eq("bucket")).
		// Limit(20).
		// Query(scyna.DB),

		scanners := scyna.DB.QueryMany("SELECT task_id FROM "+scyna_const.TODO_TABLE+" WHERE bucket = ? LIMIT 20", bucket)
		for scanners.Next() {
			var task int64
			if err := scanners.Scan(&task); err != nil {
				log.Println("scheduler.worker.execute: scan task: " + err.Error())
				return
			}
			tasks = append(tasks, task)
		}
		if len(tasks) == 0 {
			break
		}

		for _, task := range tasks {
			// qCheck: qb.Insert(scyna_const.DOING_TABLE).
			// Columns("bucket", "task_id").
			// Unique().
			// TTL(60 * time.Second).
			// Query(scyna.DB),
			if err := scyna.DB.Execute("INSERT INTO "+scyna_const.DOING_TABLE+
				" (bucket, task_id) VALUES (?, ?) IF NOT EXISTS USING TTL 60", bucket, task); err == nil {
				// if applied, _ := w.qCheck.Bind(bucket, task).ExecCAS(); applied {
				w.process(bucket, task)
			}
		}
	}
}

func (w *worker) process(bucket int64, id int64) {
	var t task
	// qGet: qb.Select(scyna_const.TASK_TABLE).
	// Columns("id", "topic", "data", "next", "interval", "loop_index", "loop_count", "done").
	// Where(qb.Eq("id")).
	// Limit(1).
	// Query(scyna.DB),

	if err := scyna.DB.QueryOne("SELECT id, topic, data, next, interval, loop_index, loop_count, done FROM "+scyna_const.TASK_TABLE+
		" WHERE id = ? LIMIT 1", id).Scan(&t.ID, &t.Topic, &t.Data, &t.Interval, &t.LoopIndex, &t.LoopCount, &t.Done); err != nil {
		// if err := w.qGet.Bind(id).Get(&t); err != nil {
		log.Print("Can not load task")
		return
	}
	if bucket != GetBucket(t.Next) {
		return
	}

	if t.Done {
		if err := scyna.DB.Execute("DELETE FROM "+scyna_const.DOING_TABLE+
			" WHERE bucket = ? AND task_id = ?", bucket, id); err != nil {
			// if err := qb.Delete(scyna_const.TODO_TABLE).
			// 	Where(qb.Eq("bucket"), qb.Eq("task_id")).
			// 	Query(scyna.DB).
			// 	Bind(bucket, id).
			// 	ExecRelease(); err != nil {
			scyna.Session.Error(err.Error())
		}
		return
	}

	scyna.JetStream.Publish(t.Topic, t.Data) /*activate task handler*/

	qBatch := scyna.DB.Session.NewBatch(gocql.LoggedBatch)
	qBatch.Query("DELETE FROM "+scyna_const.TODO_TABLE+" WHERE bucket = ? AND task_id = ?;", bucket, id) /* remove old task from todolist */

	t.LoopIndex++
	if t.LoopIndex < t.LoopCount {
		t.Next = t.Next.Add(time.Second * time.Duration(t.Interval)) /* calculate next */
		nextBucket := GetBucket(t.Next)
		qBatch.Query("INSERT INTO "+scyna_const.TODO_TABLE+" (bucket, task_id) VALUES (?, ?);", nextBucket, t.ID) /* add new task to todo list */
		qBatch.Query("UPDATE "+scyna_const.TASK_TABLE+" SET next = ?, loop_index = ?  WHERE id = ?;", t.Next, t.LoopIndex, t.ID)
	} else {
		qBatch.Query("UPDATE "+scyna_const.TASK_TABLE+" SET done = true WHERE id = ?;", t.ID)
	}

	if err := scyna.DB.Session.ExecuteBatch(qBatch); err != nil {
		scyna.Session.Error(err.Error())
	}
}
