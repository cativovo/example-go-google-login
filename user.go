package main

type User struct {
	ID        string `json:"id" doc:"Unique identifier" example:"123"`
	FirstName string `json:"first_name" doc:"First name" example:"Juan"`
	LastName  string `json:"last_name" doc:"Last name" example:"Usa"`
	Email     string `json:"email" doc:"Email" example:"juanusa@example.com"`
	AvatarURL string `json:"avatar_url" doc:"URL of the avatar" example:"https://www.images.com/123"`
}
