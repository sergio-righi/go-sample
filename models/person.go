package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	WorkAt               interface{}        `bson:"work_at" json:"work_at"`
	HonorificPrefix      interface{}        `bson:"honorific_prefix" json:"honorific_prefix"`
	Education            interface{}        `bson:"education" json:"education"`
	Spouse               interface{}        `bson:"spouse" json:"spouse"`
	Sport                interface{}        `bson:"sport" json:"sport"`
	Spoken               interface{}        `bson:"spoken" json:"spoken"`
	Pseudonym            string             `bson:"pseydonym" json:"pseudonym"`
	End                  string             `bson:"end" json:"end"`
	City                 string             `bson:"city" json:"city"`
	ParticipationOf      interface{}        `bson:"participation_of" json:"participation_of"`
	Known                interface{}        `bson:"known" json:"known"`
	Country              interface{}        `bson:"country" json:"country"`
	RecordLabel          interface{}        `bson:"record_label" json:"record_label"`
	SocialClassification string             `bson:"social_classification" json:"social_classification"`
	NativeName           string             `bson:"native_name" json:"native_name"`
	Employment           interface{}        `bson:"employement" json:"employment"`
	Start                string             `bson:"start" json:"start"`
	SexOrientation       interface{}        `bson:"sex_orientation" json:"sex_orientation"`
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EyeColor             interface{}        `bson:"eye_color" json:"eye_color"`
	Handness             interface{}        `bson:"handness" json:"handness"`
	Website              interface{}        `bson:"website" json:"website"`
	Height               interface{}        `bson:"height" json:"height"`
	Occupation           interface{}        `bson:"occupation" json:"occupation"`
	HairColor            interface{}        `bson:"hair_color" json:"hair_color"`
	Instrument           interface{}        `bson:"instrument" json:"instrument"`
	Residence            interface{}        `bson:"residence" json:"residence"`
	SkinColor            string             `bson:"skin_color" json:"skin_color"`
	BloodType            string             `bson:"blood_type" json:"blood_type"`
	Wears                interface{}        `bson:"wears" json:"wears"`
	Lost                 []interface{}      `bson:"lost" json:"lost"`
	NobleTitle           interface{}        `bson:"noble_title" json:"noble_title"`
	Speciality           string             `bson:"speciality" json:"speciality"`
	Instagram            interface{}        `bson:"instagram" json:"instagram"`
	PoliticalParty       string             `bson:"political_party" json:"political_party"`
	Mass                 string             `bson:"mass" json:"mass"`
	Nickname             string             `bson:"nickname" json:"nickname"`
	Children             string             `bson:"children" json:"children"`
	OwnerOf              []interface{}      `bson:"owner_of" json:"owner_of"`
	Name                 string             `bson:"name" json:"name"`
	NetWorth             interface{}        `bson:"net_worth" json:"net_worth"`
	Ethnic               interface{}        `bson:"ethnic" json:"ethnic"`
	Won                  []interface{}      `bson:"won" json:"won"`
	Field                interface{}        `bson:"field" json:"field"`
	Birthday             string             `bson:"birthday" json:"birthday"`
	Religion             interface{}        `bson:"religion" json:"religion"`
	Image                *string            `bson:"image,omitempty" json:"image"`
	ParticipantIn        interface{}        `bson:"participant_in" json:"participant_in"`
	Voice                interface{}        `bson:"voice" json:"voice"`
}
