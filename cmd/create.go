/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	pb "github.com/anmho/notectl/gen/proto/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var title string
var content string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called", content, title)
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		client := pb.NewNoteServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		note, err := client.CreateNote(c, &pb.CreateNoteRequest{
			Title:   title,
			Content: content,
		})
		if err != nil {
			panic(err)
		}

		log.Println(note)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	createCmd.Flags().StringVarP(&content, "content", "c", "", "note content")
	createCmd.Flags().StringVarP(&title, "title", "t", "", "title of the notes item")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


syntax= "proto3";
package notes;

option go_package = "./notes";


import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";

message Note {
string id = 1;
string title = 2;
string content = 3;
google.protobuf.Timestamp createdAt = 4;
google.protobuf.Timestamp updatedAt = 5;
}

message CreateNoteRequest {
string title = 1;
string content = 2;
}

message GetNoteRequest {
string id = 1;
}

message NoteList {
repeated Note notes = 1;
}

message DeleteNoteRequest {
string id = 1;
}

service NoteService {
rpc CreateNote(CreateNoteRequest) returns (Note) {}
rpc GetNote(GetNoteRequest) returns (Note) {}
rpc ListNotes(google.protobuf.Empty) returns (NoteList) {}
rpc DeleteNote(DeleteNoteRequest) returns (google.protobuf.Empty) {}
}


