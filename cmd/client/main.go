package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"go-grpc-rest-demo/internal/client"

	"github.com/spf13/cobra"
)

var (
	clientConfig *client.Config
	cli          client.Client
)

func main() {
	clientConfig = client.DefaultConfig()

	rootCmd := &cobra.Command{
		Use:   "client",
		Short: "CLI client for the gRPC REST demo",
		Long:  "A command line interface to interact with the user and product services",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cli, err = client.NewClient(clientConfig)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cli != nil {
				_ = cli.Close()
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&clientConfig.Mode, "mode", "m", clientConfig.Mode, "Client mode: grpc, rest")
	rootCmd.PersistentFlags().StringVar(&clientConfig.GRPCAddr, "grpc-addr", clientConfig.GRPCAddr, "gRPC server address")
	rootCmd.PersistentFlags().StringVar(&clientConfig.RESTAddr, "rest-addr", clientConfig.RESTAddr, "REST server address")
	rootCmd.PersistentFlags().DurationVar(&clientConfig.Timeout, "timeout", clientConfig.Timeout, "Request timeout")
	rootCmd.PersistentFlags().StringVar(&clientConfig.OutputFormat, "output", clientConfig.OutputFormat, "Output format: json, table")
	rootCmd.PersistentFlags().BoolVarP(&clientConfig.Verbose, "verbose", "v", clientConfig.Verbose, "Verbose output")

	rootCmd.AddCommand(userCommands(), productCommands())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func userCommands() *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management commands",
		Long:  "Commands to manage users (create, get, update, delete, list)",
	}

	createUserCmd := &cobra.Command{
		Use:   "create [username] [email] [full_name]",
		Short: "Create a new user",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var result any
			var err error

			if clientConfig.Mode == "grpc" {
				result, err = cli.CreateUserGRPC(cmd.Context(), args[0], args[1], args[2])
			} else {
				result, err = cli.CreateUserREST(cmd.Context(), args[0], args[1], args[2])
			}
			printResult(result, err, "create user")
		},
	}

	getUserCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a user by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var result any
			var err error

			if clientConfig.Mode == "grpc" {
				result, err = cli.GetUserGRPC(cmd.Context(), args[0])
			} else {
				result, err = cli.GetUserREST(cmd.Context(), args[0])
			}
			printResult(result, err, "get user")
		},
	}

	deleteUserCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cli.DeleteUser(cmd.Context(), args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete user: %v\n", err)
				return
			}
			fmt.Printf("User %s deleted successfully\n", args[0])
		},
	}

	userCmd.AddCommand(createUserCmd, getUserCmd, deleteUserCmd)
	return userCmd
}

func productCommands() *cobra.Command {
	productCmd := &cobra.Command{
		Use:   "product",
		Short: "Product management commands",
		Long:  "Commands to manage products (create, get, search)",
	}

	createProductCmd := &cobra.Command{
		Use:   "create [name] [description] [price] [quantity] [category]",
		Short: "Create a new product",
		Args:  cobra.ExactArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			price, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid price: %v\n", err)
				return
			}
			quantity, err := strconv.ParseInt(args[3], 10, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid quantity: %v\n", err)
				return
			}

			var result any
			if clientConfig.Mode == "grpc" {
				result, err = cli.CreateProductGRPC(cmd.Context(), args[0], args[1], args[4], price, int32(quantity))
			} else {
				result, err = cli.CreateProductREST(cmd.Context(), args[0], args[1], args[4], price, int32(quantity))
			}
			printResult(result, err, "create product")
		},
	}

	getProductCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a product by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var result any
			var err error

			if clientConfig.Mode == "grpc" {
				result, err = cli.GetProductGRPC(cmd.Context(), args[0])
			} else {
				result, err = cli.GetProductREST(cmd.Context(), args[0])
			}
			printResult(result, err, "get product")
		},
	}

	productCmd.AddCommand(createProductCmd, getProductCmd)
	return productCmd
}

func printResult(v any, err error, operation string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to %s (%s): %v\n", operation, clientConfig.Mode, err)
		return
	}
	printJSON(v)
}

func printJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}
