package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	"github.com/scyna/engine/gateway"
	"github.com/scyna/engine/manager/generator"
	"github.com/scyna/engine/manager/logging"
	"github.com/scyna/engine/manager/scheduler"
	"github.com/scyna/engine/manager/session"
	"github.com/scyna/engine/manager/setting"
	"github.com/scyna/engine/manager/trace"
	"github.com/scyna/engine/proxy"
)

const MODULE_CODE = "scyna_engine"

func main() {
	managerPort := flag.String("manager_port", "8081", "Manager Port")
	proxyPort := flag.String("proxy_port", "8080", "Proxy Port")
	gatewayPort := flag.String("gateway_port", "8443", "GateWay Port")

	natsUrl := flag.String("nats_url", "127.0.0.1", "Nats URL")
	natsUsername := flag.String("nats_username", "", "Nats Username")
	natsPassword := flag.String("nats_password", "", "Nats Password")
	dbHost := flag.String("db_host", "127.0.0.1", "DB Host")
	dbUsername := flag.String("db_username", "", "DB Username")
	dbPassword := flag.String("db_password", "", "DB Password")
	dbLocation := flag.String("db_location", "", "DB Location")
	secret := flag.String("secret", "123456", "scyna Manager Secret")

	certificateEnable := flag.Bool("certificateEnable", false, "Certificate Key")
	certificateFile := flag.String("certificateFile", "", "Certificate Key")
	certificateKey := flag.String("certificateKey", "", "Certificate File")

	flag.Parse()
	config := scyna_proto.Configuration{
		NatsUrl:      *natsUrl,
		NatsUsername: *natsUsername,
		NatsPassword: *natsPassword,
		DBHost:       *dbHost,
		DBUsername:   *dbUsername,
		DBPassword:   *dbPassword,
		DBLocation:   *dbLocation,
	}

	/* Init module */
	scyna.DirectInit(MODULE_CODE, &config)
	defer scyna.Release()
	generator.Init()
	session.Init(MODULE_CODE, *secret)
	scyna.UseDirectLog(5)

	/* generator */
	scyna.RegisterEndpoint(scyna_const.GEN_GET_ID_URL, generator.GetID)
	scyna.RegisterEndpoint(scyna_const.GEN_GET_SN_URL, generator.GetSN)

	/*logging*/
	scyna.RegisterSignal(scyna_const.LOG_CREATED_CHANNEL, logging.Write)

	/*trace*/
	scyna.RegisterSignal(scyna_const.TRACE_CREATED_CHANNEL, trace.TraceCreated)
	scyna.RegisterSignal(scyna_const.TAG_CREATED_CHANNEL, trace.TagCreated)
	scyna.RegisterSignal(scyna_const.ENDPOINT_DONE_CHANNEL, trace.ServiceDone)

	/*setting*/
	scyna.RegisterEndpoint(scyna_const.SETTING_READ_URL, setting.Read)
	scyna.RegisterEndpoint(scyna_const.SETTING_WRITE_URL, setting.Write)
	scyna.RegisterEndpoint(scyna_const.SETTING_REMOVE_URL, setting.Remove)

	/* task */
	scyna.RegisterEndpoint(scyna_const.START_TASK_URL, scheduler.StartTask)
	scyna.RegisterEndpoint(scyna_const.STOP_TASK_URL, scheduler.StopTask)

	/* Update config */
	setting.UpdateDefaultConfig(&config)

	const DEFAULT_CERT_FILE = ".cert/localhost.crt"
	const DEFAULT_CERT_KEY = ".cert/localhost.key"

	if *certificateEnable && (*certificateFile == "" || *certificateKey == "") {
		*certificateFile = DEFAULT_CERT_FILE
		*certificateKey = DEFAULT_CERT_KEY
	}

	go func() {
		gateway_ := gateway.NewGateway()
		log.Println("Scyna Gateway Start with port " + *gatewayPort)
		scyna.RegisterEndpoint(gateway.ADD_PUBLIC_ENDPOINT_URL, gateway.AddPublicEndpoint)
		scyna.RegisterEndpoint(gateway.REMOVE_PUBLIC_ENDPOINT_URL, gateway.RemovePublicEndpoint)

		if *certificateEnable {
			if err := http.ListenAndServeTLS(":"+*gatewayPort, *certificateFile, *certificateKey, gateway_); err != nil {
				log.Println("Gateway: " + err.Error())
			}
		} else {
			if err := http.ListenAndServe(":"+*gatewayPort, gateway_); err != nil {
				log.Println("Gateway: " + err.Error())
			}
		}
	}()

	go func() {
		proxy_ := proxy.NewProxy()
		log.Println("Scyna Proxy Start with port " + *proxyPort)

		if *certificateEnable && *certificateFile != "" {
			if err := http.ListenAndServeTLS(":"+*proxyPort, *certificateFile, *certificateKey, proxy_); err != nil {
				log.Println("Proxy: " + err.Error())
			}
		} else {
			if err := http.ListenAndServe(":"+*proxyPort, proxy_); err != nil {
				log.Println("Proxy: " + err.Error())
			}
		}
	}()

	/* Start worker */
	scheduler.Start(time.Second * 10)
	/*session*/
	scyna.RegisterSignal(scyna_const.SESSION_END_CHANNEL, session.End)
	scyna.RegisterSignal(scyna_const.SESSION_UPDATE_CHANNEL, session.Update)
	http.HandleFunc(scyna_const.SESSION_CREATE_URL, session.Create)
	log.Println("Scyna Manager Start with port " + *managerPort)
	if *certificateEnable && *certificateFile != "" {
		if err := http.ListenAndServeTLS(":"+*managerPort, *certificateFile, *certificateKey, nil); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(":"+*managerPort, nil); err != nil {
			panic(err)
		}
	}
}
