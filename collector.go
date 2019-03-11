package main

import (
      "time"
      "bytes"
      "sync"
      "strconv"
      "io/ioutil"
      "encoding/json"
      "net"
      "net/http"
      "github.com/prometheus/client_golang/prometheus"
    _ "github.com/patrickmn/go-cache"
  log "github.com/Sirupsen/logrus"
      "github.com/pkg/errors"
      "database/sql"
    _ "gopkg.in/goracle.v2"
)

var humanTaskMismatchesForNotificationsMetric = prometheus.NewGaugeVec(
  prometheus.GaugeOpts{
    Name: "business_ext_notifications_human_tasks_mismatch_count",
    Help: "Number of notification processes that do not have the minimum number of created human tasks",
  },
  []string{"country_abbreviation","notification_type"},
)

var netTransport = &http.Transport{
  Dial: (&net.Dialer{
    Timeout: 5 * time.Second,
  }).Dial,
  TLSHandshakeTimeout: 5 * time.Second,
}

var netClient = &http.Client{
  Timeout: time.Second * 10,
  Transport: netTransport,
}


// DB CONNECTION ///////////////////////////////////////////////////////////////
var oracleDsn = "oracle://?sysdba=1&prelim=1"
var db        *sql.DB //= sql.Open("goracle", oracleDsn)

// NET CONNECTION //////////////////////////////////////////////////////////////
var mwsUrl          = "http://localhost:8888/"
var mwsCredUsername = "Administrator"
var mwsCredPassword = "manage"
var jsonContentType = "application/json"




func init() {
  log.Info("INIT RUNNING")
  var err error
  db, err = sql.Open("goracle", oracleDsn)
  if err != nil {
      log.Info("WUMS")
    log.Fatal(errors.Wrap(err, oracleDsn))
  }
  //log.Info("DEFER CLOSES")
  //defer db.Close()
}


















//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type DmsCollector struct {
	missingProcessesCount    *prometheus.Desc
	missingGatewayTasksCount *prometheus.Desc
	missingFirstTasksCount   *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newDmsCollector() *DmsCollector {
	return &DmsCollector{
		missingProcessesCount: prometheus.NewDesc("missing_processes_count",
			"Number of missing receiver country processes.",
			nil, nil,
		),
  	missingGatewayTasksCount: prometheus.NewDesc("missing_gateway_tasks_count",
  		"Number of missing gateway tasks in the receiver country processes.",
  		nil, nil,
  	),
    missingFirstTasksCount: prometheus.NewDesc("missing_first_tasks_count",
    	"Number of missing first tasks in the receiver country processes.",
    	nil, nil,
    ),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *DmsCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.missingProcessesCount
	ch <- collector.missingGatewayTasksCount
	ch <- collector.missingFirstTasksCount
}

//Collect implements required collect function for all promehteus collectors
func (collector *DmsCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
//TODO: DO ALL THE STUFF TO DO!
countRetrievalProcessesExpectedPerSenderProcess()

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(collector.missingProcessesCount,    prometheus.GaugeValue, 42)
	ch <- prometheus.MustNewConstMetric(collector.missingGatewayTasksCount, prometheus.GaugeValue, 43)
	ch <- prometheus.MustNewConstMetric(collector.missingFirstTasksCount,   prometheus.GaugeValue, 44)

}


















// TODO: https://better-coding.com/solved-njs-045-dpi-1047-64-bit-oracle-client-library-cannot-be-loaded/

var countRetrievalProcessesExpectedPerSenderProcessLatestUpdate = 0
var countRetrievalProcessesExpectedPerSenderProcessMap      = make(map[int]int)
var countRetrievalProcessesExpectedPerSenderProcessMapMutex = &sync.Mutex{}

// count all theoretically created target processes
func countRetrievalProcessesExpectedPerSenderProcess() {

  var selectStatement = `
    select id, name from users where id = ? and updatetimestamp > ?
  `
  // read from db
  rows, err := db.Query(selectStatement, 1, countRetrievalProcessesExpectedPerSenderProcessLatestUpdate)
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()

  // read result
  var senderNotificationId         int
  var expectedReceiverProcessCount int
  for rows.Next() {
    err := rows.Scan(&senderNotificationId, &expectedReceiverProcessCount)
    if err != nil {
      log.Fatal(err)
    }

    // handle row
    // set latest update date
    //countRetrievalProcessesExpectedPerSenderProcessLatestUpdate

    log.Info("TXID: " + strconv.Itoa(senderNotificationId) + ": " + strconv.Itoa(expectedReceiverProcessCount))
    countRetrievalProcessesExpectedPerSenderProcessMapMutex.Lock()
    countRetrievalProcessesExpectedPerSenderProcessMap[senderNotificationId] = expectedReceiverProcessCount
    countRetrievalProcessesExpectedPerSenderProcessMapMutex.Unlock()
  }
  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }

}
/*
// count really created target processes
func countRetrievalProcessesExistingPerSenderProcess() map[int]int {
}

// count theoretically created gateway tasks
func countGatewayTasksExpectedPerRetrieverProcess() map[int]int {
}

// count really created gateway tasks
func countGatewayTasksExistingPerRetrieverProcess() map[int]int {
}
*/
// count theoretically created 1ststep tasks
func count1ststepTasksExpectedPerRetrieverProcess() {

  var url = mwsUrl + "services/bizPolicy/task/searchTasks"

  var jsonRequestContent = []byte(`
    {
      "includeTaskData" : true,
      "taskSearchQuery" : {
        "terms" : [
          "x=y"
        ]
      }
    }
  `)

  request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequestContent))
  request.Header.Set("Content-Type", jsonContentType)
  request.SetBasicAuth(mwsCredUsername, mwsCredPassword)

  response, err := netClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

  log.Info(response)

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    log.Fatalln(err)
  }

  var responseData string
  json.Unmarshal(body, &responseData)
}
/*
// count really created 1ststep tasks
func count1ststepTasksExistingPerRetrieverProcess() map[int]int {
}
*/


func init_old() {
log.Info("000")
  findTasks()
  return


  prometheus.MustRegister(humanTaskMismatchesForNotificationsMetric)

  // initially search for all sent notifications and add them to the cache
  // TODO: :-)

  // initially search for all retrieved notifications and add the number of created tasks to the cache
  // TODO: :-)

/*
  // get number of theoretical human tasks for each sender country process and store in memory
  // if possible filter done processes before - or filter by creation time (1 month)
  // TODO: CACHE IS NOT ALLOWED TO EXPIRE!!!
  expectedCount := cache.New(5*time.Minute, 10*time.Minute)

  // set falue for sent notification
  expectedCount.Set("57458", 42, cache.NoExpiration)

  // Get the string associated with the key "foo" from the cache
	expectedCount, found := expectedCount.Get("foo")
	if !found {
		// calculate
	}
*/

  // periodically search for newly sent notifications and add them to the cache
  // TODO: :-)

  // periodically search for newly retrieved notifications and add the number of created tasks to the cache
  // TODO: :-)

/*
  go func() {
    for {
      // https://github.com/go-goracle/goracle

      // TODO: CHECK THIS OUT!
      // db.QueryContext(goracle.ContextWithLog(ctx, logger.Log), qry)

		  // connect to database
		  // TODO: Could the connection stay open and only be reestablished if closed by server?
		  //       http://go-database-sql.org/accessing.html
		  dsn := "oracle://?sysdba=1&prelim=1"
      var db, err = sql.Open("goracle", dsn)

		  // do some errorhandling
      if err != nil {
        log.Fatal(errors.Wrap(err, dsn))
      }

		  // ensure that the database is closed at the end
      defer db.Close()

		  // example for multiple calls with different content - will not be used
		  // db.Exec("INSERT INTO table (a, b) VALUES (:1, :2)", []int{1, 2}, []string{"a", "b"})

		  rows, err := db.Query("select id, name from users where id = ?", 1)
      if err != nil {
        log.Fatal(err)
      }
      defer rows.Close()

      for rows.Next() {
        err := rows.Scan(&id, &name)
        if err != nil {
          log.Fatal(err)
        }
        log.Info(id + name)
      }
      err = rows.Err()
      if err != nil {
        log.Fatal(err)
      }

      humanTaskMismatchesForNotificationsMetric.WithLabelValues("DE","Ware aus dem Verkauf").Inc() //Add(42)//Set(123)

			time.Sleep(time.Duration(15) * time.Second)
    }
  }()
*/
  go func() {
    for {
      //https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
      //https://github.com/mwitkow/go-http-dialer
      url := "http://www.google.de/"
/*
      var netTransport = &http.Transport{
        Dial: (&net.Dialer{
          Timeout: 5 * time.Second,
        }).Dial,
        TLSHandshakeTimeout: 5 * time.Second,
      }

      var netClient = &http.Client{
        Timeout: time.Second * 10,
        Transport: netTransport,
      }
*/
      response, _ := netClient.Get(url)

      log.Info(response)

      humanTaskMismatchesForNotificationsMetric.WithLabelValues("DE","Ware aus dem Verkauf").Inc() //Add(42)//Set(123)

			time.Sleep(time.Duration(15) * time.Second)
		}
  }()

}


func findTasks() int {
log.Info("0")
  //url := "https://mwslb/services/bizPolicy/task/searchTasks"
  url := "http://localhost:8888/services/bizPolicy/task/searchTasks"
  contentType := "application/json"


  var jsonStr = []byte(`
    {
      "includeTaskData" : true,
      "taskSearchQuery" : {
        "terms" : [
          "x=y"
        ]
      }
    }
  `)
/*
	jsonData := map[string]interface{}{
		"includeTaskData": true,
		"taskSearchQuery": map[string]string{
			"terms": string[]{
        "x=y"
      },
		},
	}

	jsonDataBytesRepresentation, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalln(err)
	}
*/
//request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonDataBytesRepresentation))
  request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  request.Header.Set("Content-Type", contentType)
  request.SetBasicAuth("Administrator", "manage")
log.Info("A")
  response, err := netClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
log.Info("B")
  log.Info(response)

  return 42
}
