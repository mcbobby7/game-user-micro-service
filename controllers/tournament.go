package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"

	"time"

	"github.com/Gameware/database"
	helper "github.com/Gameware/helpers"
	"github.com/Gameware/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var tournamentCollection *mongo.Collection = database.OpenCollection(database.Client, "tournament")

func SaveTournament() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var tournament models.Tournament

		if err := c.BindJSON(&tournament); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "hasError": true})
			return
		}

		validationErr := validate.Struct(tournament)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error(), "hasError": true})
			return
		}

		tournament.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tournament.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tournament.ID = primitive.NewObjectID()
		tournament.TournamentId = tournament.ID.Hex()
		tournament.User_id = c.GetString("uid")

		resultInsertionNumber, insertErr := tournamentCollection.InsertOne(ctx, tournament)
		if insertErr !=nil {
			msg := fmt.Sprintf("item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg, "hasError": true})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"message": "request processed successfullt", "data":tournament, "hasError": false, "insertId": resultInsertionNumber})
	}
}
func GetTournament() gin.HandlerFunc{
	return func(c *gin.Context){
		id := c.Param("id")
		// log.Fatal(id)

		if err := helper.MatchUserTypeToUid(c, id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error(), "hasError": true})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var tournament models.Tournament
		err := tournamentCollection.FindOne(ctx, bson.M{"tournamentid":id}).Decode(&tournament)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "hasError": true})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "request processed successfullt", "tournament":tournament, "hasError": false})
	}
}
func GetTournaments() gin.HandlerFunc{
	return func(c *gin.Context){
		// if err := helper.CheckUserType(c, "ADMIN"); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		// 	return
		// }
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <1{
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 !=nil || page<1{
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}
		result,err := tournamentCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items", "hasError": true})
			return
		}
		var data []bson.M
		if err = result.All(ctx, &data); err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "request processed successfullt", "tournaments":data[0], "hasError": false})}
}