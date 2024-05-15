/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	pb "github.com/anmho/notectl/gen/proto/notes"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"text/tabwriter"
	"time"
)

var id string
var isAll bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <note_id>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !isAll || len(args) > 0 && isAll {
			err := cmd.Usage()
			if err != nil {
				panic(err)
			}
			return
		}

		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}

		client := pb.NewNoteServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if isAll {
			notes, err := client.ListNotes(c, &emptypb.Empty{})
			if err != nil {
				panic(err)
			}

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '-', tabwriter.AlignRight|tabwriter.Debug)
			fmt.Fprintln(w, "id\ttitle\tcontent\tcreatedAt\tupdatedAt\t")

			for _, note := range notes.Notes {

				//fmt.Println(note)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n",
					note.Id,
					note.Title,
					note.Content,
					note.CreatedAt.AsTime().Format(time.RFC822),
					note.UpdatedAt.AsTime().Format(time.RFC822),
				)
			}
			w.Flush()

		} else {
			id = args[0]
			req := &pb.GetNoteRequest{
				Id: id,
			}

			note, err := client.GetNote(c, req)
			if err != nil {
				panic(err)
			}

			fmt.Println(note)
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.
	getCmd.Flags().BoolVarP(&isAll, "all", "a", false, "notectl create -a")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
