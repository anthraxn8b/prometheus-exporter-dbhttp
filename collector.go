package main

import (
      "time"
      "net"
      "net/http"
      "github.com/prometheus/client_golang/prometheus"
  log "github.com/Sirupsen/logrus"
    _ "github.com/pkg/errors"
    _ "database/sql"
    _ "gopkg.in/goracle.v2"
)

var humanTaskMismatchesForNotificationsMetric = prometheus.NewGaugeVec(
  prometheus.GaugeOpts{
    Name: "business_ext_notifications_human_tasks_mismatch_count",
    Help: "Number of notification processes that do not have the minimum number of created human tasks",
  },
  []string{"country_abbreviation","notification_type"},
)

func init() {
  prometheus.MustRegister(humanTaskMismatchesForNotificationsMetric)
   
  // periodically fetch the database information
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

      response, _ := netClient.Get(url)
   
      log.Info(response)
      
      humanTaskMismatchesForNotificationsMetric.WithLabelValues("DE","Ware aus dem Verkauf").Inc() //Add(42)//Set(123)

			time.Sleep(time.Duration(15) * time.Second)
		}
  }()

}
