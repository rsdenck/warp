package contacts

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/output"
	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ContactsCmd = &cobra.Command{
		Use:   "contacts",
		Short: "Contacts operations",
		Long:  `Commands for IceWarp Contacts API`,
	}
)

var listContactsCmd = &cobra.Command{
	Use:   "list",
	Short: "List contacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewContactsClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		folderID, _ := cmd.Flags().GetString("folder")

		contacts, err := client.ListContacts(folderID)
		if err != nil {
			return fmt.Errorf("failed to list contacts: %w", err)
		}

		t := output.NewTable("CONTACTS")
		t.AppendHeader(table.Row{"First Name", "Last Name", "Email", "ID"})
		
		for _, c := range contacts {
			t.AppendRow(table.Row{c.FirstName, c.LastName, c.Email, c.ID})
		}
		
		t.Render()
		return nil
	},
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List contact groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewContactsClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		groups, err := client.ListGroups()
		if err != nil {
			return fmt.Errorf("failed to list groups: %w", err)
		}

		t := output.NewTable("CONTACT GROUPS")
		t.AppendHeader(table.Row{"Name", "ID"})
		
		for _, g := range groups {
			t.AppendRow(table.Row{g.Name, g.ID})
		}
		
		t.Render()
		return nil
	},
}

var createContactCmd = &cobra.Command{
	Use:   "create [email]",
	Short: "Create a new contact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewContactsClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		firstName, _ := cmd.Flags().GetString("first-name")
		lastName, _ := cmd.Flags().GetString("last-name")
		phone, _ := cmd.Flags().GetString("phone")
		mobile, _ := cmd.Flags().GetString("mobile")
		company, _ := cmd.Flags().GetString("company")
		jobTitle, _ := cmd.Flags().GetString("job")
		address, _ := cmd.Flags().GetString("address")
		city, _ := cmd.Flags().GetString("city")
		state, _ := cmd.Flags().GetString("state")
		country, _ := cmd.Flags().GetString("country")
		zipCode, _ := cmd.Flags().GetString("zip")
		webSite, _ := cmd.Flags().GetString("website")
		notes, _ := cmd.Flags().GetString("notes")
		folderID, _ := cmd.Flags().GetString("folder")

		contact := &sdk.Contact{
			Email:     args[0],
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phone,
			Mobile:    mobile,
			Company:   company,
			JobTitle:  jobTitle,
			Address:   address,
			City:      city,
			State:     state,
			Country:   country,
			ZipCode:   zipCode,
			WebSite:   webSite,
			Notes:     notes,
			FolderID:  folderID,
		}

		result, err := client.CreateContact(contact)
		if err != nil {
			return fmt.Errorf("failed to create contact: %w", err)
		}

		fmt.Printf("Contact '%s %s' created (ID: %s)\n", result.FirstName, result.LastName, result.ID)
		return nil
	},
}

var updateContactCmd = &cobra.Command{
	Use:   "update [contact-id]",
	Short: "Update a contact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewContactsClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		firstName, _ := cmd.Flags().GetString("first-name")
		lastName, _ := cmd.Flags().GetString("last-name")
		email, _ := cmd.Flags().GetString("email")
		phone, _ := cmd.Flags().GetString("phone")
		mobile, _ := cmd.Flags().GetString("mobile")
		company, _ := cmd.Flags().GetString("company")
		jobTitle, _ := cmd.Flags().GetString("job")

		contact := &sdk.Contact{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phone,
			Mobile:    mobile,
			Company:   company,
			JobTitle:  jobTitle,
		}

		result, err := client.UpdateContact(args[0], contact)
		if err != nil {
			return fmt.Errorf("failed to update contact: %w", err)
		}

		fmt.Printf("Contact updated (ID: %s)\n", result.ID)
		return nil
	},
}

var deleteContactCmd = &cobra.Command{
	Use:   "delete [contact-id]",
	Short: "Delete a contact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewContactsClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteContact(args[0]); err != nil {
			return fmt.Errorf("failed to delete contact: %w", err)
		}

		fmt.Printf("Contact deleted successfully\n")
		return nil
	},
}

var createGroupCmd = &cobra.Command{
	Use:   "create-group [name]",
	Short: "Create a new contact group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewContactsClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		parentID, _ := cmd.Flags().GetString("parent")

		group, err := client.CreateGroup(args[0], parentID)
		if err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		fmt.Printf("Group '%s' created (ID: %s)\n", group.Name, group.ID)
		return nil
	},
}

var deleteGroupCmd = &cobra.Command{
	Use:   "delete-group [group-id]",
	Short: "Delete a contact group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewContactsClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteGroup(args[0]); err != nil {
			return fmt.Errorf("failed to delete group: %w", err)
		}

		fmt.Printf("Group deleted successfully\n")
		return nil
	},
}

func init() {
	ContactsCmd.AddCommand(listContactsCmd)
	ContactsCmd.AddCommand(listGroupsCmd)
	ContactsCmd.AddCommand(createContactCmd)
	ContactsCmd.AddCommand(updateContactCmd)
	ContactsCmd.AddCommand(deleteContactCmd)
	ContactsCmd.AddCommand(createGroupCmd)
	ContactsCmd.AddCommand(deleteGroupCmd)

	listContactsCmd.Flags().StringP("folder", "f", "", "Folder ID")

	createContactCmd.Flags().StringP("first-name", "n", "", "First name")
	createContactCmd.Flags().StringP("last-name", "l", "", "Last name")
	createContactCmd.Flags().StringP("phone", "p", "", "Phone")
	createContactCmd.Flags().StringP("mobile", "m", "", "Mobile")
	createContactCmd.Flags().StringP("company", "c", "", "Company")
	createContactCmd.Flags().StringP("job", "j", "", "Job title")
	createContactCmd.Flags().StringP("address", "a", "", "Address")
	createContactCmd.Flags().StringP("city", "i", "", "City")
	createContactCmd.Flags().StringP("state", "s", "", "State")
	createContactCmd.Flags().StringP("country", "o", "", "Country")
	createContactCmd.Flags().StringP("zip", "z", "", "Zip code")
	createContactCmd.Flags().StringP("website", "w", "", "Website")
	createContactCmd.Flags().StringP("notes", "t", "", "Notes")
	createContactCmd.Flags().StringP("folder", "f", "", "Folder ID")

	updateContactCmd.Flags().StringP("first-name", "n", "", "First name")
	updateContactCmd.Flags().StringP("last-name", "l", "", "Last name")
	updateContactCmd.Flags().StringP("email", "e", "", "Email")
	updateContactCmd.Flags().StringP("phone", "p", "", "Phone")
	updateContactCmd.Flags().StringP("mobile", "m", "", "Mobile")
	updateContactCmd.Flags().StringP("company", "c", "", "Company")
	updateContactCmd.Flags().StringP("job", "j", "", "Job title")

	createGroupCmd.Flags().StringP("parent", "p", "", "Parent group ID")
}
