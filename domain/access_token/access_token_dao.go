package access_token

import (
	"github.com/StarsPoker/loginBackend/logger"
	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/mongo_utils"

	"github.com/StarsPoker/loginBackend/datasources/mongo/stars_mongo"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	database   = "stars"
	collection = "access_token"
	records_collection = "access_token_records"
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

func CreateRecord(accessToken AccessToken) *rest_errors.RestErr {
	session, _ := stars_mongo.GetSession()
	defer session.Close()
	
	col := session.DB(database).C(records_collection)

	err := col.Insert(&accessToken)
	if err != nil {
		logger.Error("error when trying to save access_token_record", err)
		return rest_errors.NewInternalServerError("database error")
	}
	
	return nil
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

func Delete(accessTokenId string) *rest_errors.RestErr {
	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)
	err := col.Remove(bson.M{"access_token": accessTokenId})

	if err != nil {
		logger.Error("error when trying to delete a access_token", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func UpdateLastInteraction(accessTokenId string) *rest_errors.RestErr {

	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)

	err := col.Update(bson.M{"access_token": accessTokenId}, bson.M{"$set": bson.M{"last_interaction": date_utils.GetNow().Unix()}})

	if err != nil {
		logger.Error("error when trying to update last interaction access token", err)
		return mongo_utils.ParseError(err)
	}

	return nil
}

func DeleteExpiredAccesTokens() *rest_errors.RestErr {
	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)

	err := col.Remove(bson.M{"expires": bson.M{"$lte": date_utils.GetNow().Unix()}})

	if err != nil {
		logger.Error("error when trying to delete a access_token", err)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}
