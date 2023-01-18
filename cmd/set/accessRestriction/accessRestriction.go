package accessRestriction

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/spf13/cobra"
)



type domain struct {
	BlacklistDomain    string `json:"blacklistdomain"`
	WhitelistDomain    string `json:"whitelistdomain"`
	DomainMod string `json:"domainmod"`
}


func NewaccessRestrictionCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Updates whitelisted/blacklisted domain/emails",
		Long:  `Use this command to update the whitelisted/blacklisted domains/emails.`,
		Example: heredoc.Doc(`$ lr set access-restriction --blacklist-domain <old-domain> --new-domain <new-domain>
		<Type> domains/emails have been updated
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.BlacklistDomain == "" && opts.WhitelistDomain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}

			if opts.DomainMod == "" {
				return &cmdutil.FlagError{Err: errors.New("`new-domain` is required argument")}
			}

			resp, err := api.GetEmailWhiteListBlackList()
			if err != nil {
				return err
			}
			if (resp.ListType == "WhiteList" && opts.BlacklistDomain != "") || (resp.ListType == "BlackList" && opts.WhitelistDomain != "") {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email Not Found. As " + resp.ListType + " Restriction Type is selected. You can change it via `lr add access-restriction`" )}
			} 
			var domain string
			if opts.BlacklistDomain != "" {
				domain = opts.BlacklistDomain
			} else if opts.WhitelistDomain != "" {
				domain = opts.WhitelistDomain
			}

			i, found := cmdutil.Find(resp.Domains, domain)
			if !found {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email not found")}
			}
			if !cmdutil.AccessRestrictionDomain.MatchString(opts.DomainMod)  {
				return &cmdutil.FlagError{Err: errors.New("Domain/Email field is invalid")}
			}
			
			_, found = cmdutil.Find(resp.Domains, opts.DomainMod)
			if found {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email is already added")}
			}
			var newDomains []string
			newDomains = resp.Domains
			newDomains[i] = opts.DomainMod
			set(resp.ListType, newDomains)
			return nil

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.BlacklistDomain, "blacklist-domain", "b", "", "Enter Old Blacklist Domain/Email Value")
	fl.StringVarP(&opts.WhitelistDomain, "whitelist-domain", "w", "", "Enter Old Whitelist Domain/Email Value")
	fl.StringVarP(&opts.DomainMod, "new-domain", "n", "", "Enter New Domain Value")

	return cmd
}

func set(listType string, domain []string) error {
	var restrictType api.RegistrationRestrictionTypeSchema
	restrictType.SelectedRestrictionType = strings.ToLower(listType)
	var AddEmail api.EmailWhiteBLackListSchema
	AddEmail.Domains = domain
	err := api.AddEmailWhitelistBlacklist(restrictType, AddEmail);
	if err != nil {
		return err
	}
	fmt.Println(listType + " domains/emails have been updated" )
	return nil
}
