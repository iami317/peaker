package plugins

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
)

func ScanMongodb(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	mgUrl := fmt.Sprintf(
		"mongodb://%v:%v@%v:%v",
		s.Username,
		url.QueryEscape(s.Password),
		s.Ip,
		s.Port)
	ctx := context.Background()
	session, err := mongo.Connect(ctx, options.Client().ApplyURI(mgUrl))

	if err == nil {
		defer session.Disconnect(ctx)
		err = session.Ping(ctx, readpref.Primary())
		if err == nil {
			result.Class = WeakPass
			result.Result = true
		}
	}

	return result
}

func UnauthorizedMongdb(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	mgUrl := fmt.Sprintf("mongodb://%v:%v@%v:%v", "", url.QueryEscape(""), s.Ip, s.Port)
	ctx := context.Background()
	session, err := mongo.Connect(ctx, options.Client().ApplyURI(mgUrl))

	if err == nil {
		defer session.Disconnect(ctx)
		err = session.Ping(ctx, readpref.Primary())
		if err == nil {
			result.Class = Unauthorized
			result.Result = true
		}
	}

	return result
}
