# go-ldap-slack-syncer: utility for synchronizing slack users with LDAP users

## Requirements for running the utility:

* MySQL database
* Access to the LDAP server
* Access to the SlackAPI
* Config file 'example file *config-example.yaml*'

### Description of startup configuration flags:
* **--apply** the absence of this flag indicates the module to work in debug mode (only write about its intentions to the channel, without performing actions).
* **--revert** enables the mode of returning to the previous state
* **--date** или **-d** the flag is used only if the revert mode is enabled. Used to specify a specific date for which to make a revert. Date format **YYYY-MM-DD**
* **--enable** the flag is used only if the revert mode is enabled. Specifies that users in the revert mode should be enabled if they have been disabled.
* **--disable** the flag is used only if the revert mode is enabled. Specifies that users in the revert mode should be disabled if they have been enabled.
* **--usermail** или **-u** the flag is used only if the revert mode is enabled. If the user's email is specified, revert mode will work only for this user.
* **--count** optional flag to specify the maximum number of synchronization users. If more users are found than indicated by the flag, they will be ignored and synchronized only when the utility is restarted.
* **--lastupdate** flag for specifying the date of the last update/activity of the slack users. Date format **YYYY-MM-DD**
* **--before** synchronize users who were active before the date specified in lastupdate
* **--after** synchronize users who were active after the date specified in the last update