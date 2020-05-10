package access_token

import (
	"github.com/StarsPoker/loginBackend/logger"
	"github.com/StarsPoker/loginBackend/utils/mongo_utils"

	"github.com/StarsPoker/loginBackend/datasources/mongo/stars_mongo"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	database   = "stars"
	collection = "access_token"
)

func GetById(accessTokenId string) (*AccessToken, *rest_errors.RestErr) {

	at := AccessToken{}

	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)

	err := col.Find(bson.M{"access_token": accessTokenId}).One(&at)

	if err != nil {
		logger.Error("error when trying to get a access token", err)
		return nil, mongo_utils.ParseError(err)
	}

	return &at, nil
}

func Create(accessToken AccessToken) *rest_errors.RestErr {
	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)
	err := col.Insert(&accessToken)

	if err != nil {
		logger.Error("error when trying to creat a access_token", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func UpdateExpirationTime(accessToken *AccessToken) (*AccessToken, *rest_errors.RestErr) {

	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)

	err := col.UpdateId(accessToken.Id, bson.M{"$set": bson.M{
		"expires": accessToken.Expires,
	}})

	if err != nil {
		logger.Error("error when trying to update a access token", err)
		return nil, mongo_utils.ParseError(err)
	}

	return accessToken, nil
}
