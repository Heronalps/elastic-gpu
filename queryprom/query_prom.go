package queryprom

import (
	"context"
	"errors"
	"log"
	"time"

	clientapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/spf13/viper"
)

/*
Query the prometheus server and return vector value
parameters:
	namespace
	queryStr
returns:
	value
	err
*/
func Query(namespace string, queryStr string) (value float64, err error) {
	// Watching namespaces usage
	client, err := clientapi.NewClient(clientapi.Config{Address: "https://prometheus.nautilus.optiputer.net"})
	if err != nil {
		log.Printf("%v", err)
		return 0, err
	}

	q := prometheus.NewAPI(client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// signature: Query(ctx context.Context, query string, ts time.Time) (model.Value, api.Error)
	// curVal is of type "model.Value"
	if curVal, err := q.Query(ctx, queryStr, time.Now()); err != nil {
		log.Printf("%v", err)
	} else {
		switch {
		case curVal.Type() == model.ValVector:
			vectorVal := curVal.(model.Vector)
			// fmt.Println("===Vector====")
			// fmt.Println(vectorVal)
			for _, elem := range vectorVal {
				for _, ns := range viper.GetStringSlice("portal.gpu_exceptions") {
					if string(elem.Metric["namespace_name"]) == ns {
						return 0, errors.New("gpu exceptions")
					}
				}

				if string(elem.Metric["namespace_name"]) == namespace {
					return float64(elem.Value), nil
				}
			}
		}
	}
	return 0, errors.New("Invalid query: Query Range is too short")
}
