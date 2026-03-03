package certs

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CertsCmd = &cobra.Command{
		Use:   "certs",
		Short: "SSL Certificates operations",
		Long:  `Commands for IceWarp SSL Certificates API`,
	}
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Certificates API",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid, err := client.Authenticate(username, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Printf("Login successful!\n")
		fmt.Printf("SID: %s\n", sid)
		return nil
	},
}

var listCertsCmd = &cobra.Command{
	Use:   "list",
	Short: "List SSL certificates",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		certs, err := client.ListCertificates()
		if err != nil {
			return fmt.Errorf("failed to list certificates: %w", err)
		}

		fmt.Println("SSL Certificates:")
		for _, c := range certs {
			fmt.Printf("  %s (ID: %s) - Domain: %s - Valid: %v, Default: %v\n", c.Name, c.ID, c.Domain, c.IsValid, c.IsDefault)
		}
		return nil
	},
}

var getInfoCmd = &cobra.Command{
	Use:   "info [cert-id]",
	Short: "Get certificate information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		cert, err := client.GetCertificateInfo(args[0])
		if err != nil {
			return fmt.Errorf("failed to get info: %w", err)
		}

		fmt.Printf("Certificate: %s\n", cert.Name)
		fmt.Printf("Domain: %s\n", cert.Domain)
		fmt.Printf("Issuer: %s\n", cert.Issuer)
		fmt.Printf("Subject: %s\n", cert.Subject)
		fmt.Printf("Fingerprint: %s\n", cert.Fingerprint)
		fmt.Printf("Valid: %v\n", cert.IsValid)
		fmt.Printf("Default: %v\n", cert.IsDefault)
		return nil
	},
}

var addCertCmd = &cobra.Command{
	Use:   "add [name] [domain]",
	Short: "Add a new SSL certificate",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		certData, _ := cmd.Flags().GetString("cert")
		keyData, _ := cmd.Flags().GetString("key")
		password, _ := cmd.Flags().GetString("password")

		err := client.AddCertificate(args[0], args[1], certData, keyData, password)
		if err != nil {
			return fmt.Errorf("failed to add certificate: %w", err)
		}

		fmt.Printf("Certificate added\n")
		return nil
	},
}

var deleteCertCmd = &cobra.Command{
	Use:   "delete [cert-id]",
	Short: "Delete an SSL certificate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		if err := client.DeleteCertificate(args[0]); err != nil {
			return fmt.Errorf("failed to delete: %w", err)
		}

		fmt.Printf("Certificate deleted\n")
		return nil
	},
}

var setDefaultCmd = &cobra.Command{
	Use:   "set-default [cert-id]",
	Short: "Set a certificate as default",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		if err := client.SetDefaultCertificate(args[0]); err != nil {
			return fmt.Errorf("failed to set default: %w", err)
		}

		fmt.Printf("Certificate set as default\n")
		return nil
	},
}

var createCSRCmd = &cobra.Command{
	Use:   "create-csr [domain] [common-name]",
	Short: "Create a Certificate Signing Request",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		org, _ := cmd.Flags().GetString("org")
		orgUnit, _ := cmd.Flags().GetString("org-unit")
		city, _ := cmd.Flags().GetString("city")
		state, _ := cmd.Flags().GetString("state")
		country, _ := cmd.Flags().GetString("country")

		csr, err := client.CreateCSR(args[0], args[1], org, orgUnit, city, state, country)
		if err != nil {
			return fmt.Errorf("failed to create CSR: %w", err)
		}

		fmt.Printf("CSR:\n%s\n", csr)
		return nil
	},
}

var exportCertCmd = &cobra.Command{
	Use:   "export [cert-id]",
	Short: "Export a certificate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCertificatesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		cert, err := client.ExportCertificate(args[0])
		if err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}

		fmt.Printf("Certificate:\n%s\n", cert)
		return nil
	},
}

func init() {
	CertsCmd.AddCommand(loginCmd)
	CertsCmd.AddCommand(listCertsCmd)
	CertsCmd.AddCommand(getInfoCmd)
	CertsCmd.AddCommand(addCertCmd)
	CertsCmd.AddCommand(deleteCertCmd)
	CertsCmd.AddCommand(setDefaultCmd)
	CertsCmd.AddCommand(createCSRCmd)
	CertsCmd.AddCommand(exportCertCmd)

	addCertCmd.Flags().StringP("cert", "c", "", "Certificate data (PEM)")
	addCertCmd.Flags().StringP("key", "k", "", "Private key data")
	addCertCmd.Flags().StringP("password", "p", "", "Certificate password")

	createCSRCmd.Flags().StringP("org", "o", "", "Organization")
	createCSRCmd.Flags().StringP("org-unit", "u", "", "Organization Unit")
	createCSRCmd.Flags().StringP("city", "i", "", "City")
	createCSRCmd.Flags().StringP("state", "s", "", "State")
	createCSRCmd.Flags().StringP("country", "t", "", "Country")
}
