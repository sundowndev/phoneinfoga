v1.0.54
----------
 * update metadata for v8.11.0

v1.0.53
----------
 * Metadata update for upstream v8.10.23

v1.0.52
----------
 * Reset italian leading zero when false, fixed bug when phonenumber struct is reused

v1.0.51
----------
 * Update metadata to upstream 8.10.21

v1.0.50
----------
 * Fix formatting of country code in out-of-country format (thanks @janh)
 * Fix FormatInOriginalFormat for numbers with national prefix (thanks @janh)
 * Fix panic due to calling proto.Merge on nil destination (thanks @janh)

v1.0.49
----------
 * fix Makefile for phoneserver

v1.0.48
----------
 * another test travis rev, ignore

v1.0.47
----------
 * test tag for travis deploy

v1.0.46
----------
 * update metadata for v8.10.19
 * remove aws-lambda-go as dependency (thanks @shaxbee)

v1.0.45
----------
 * Update metadata to fix Mexican formatting (thanks @bvisness)
 * Add tests specifically for Mexico (thanks @bvisness)

v1.0.44
----------
 * update metadata for v8.10.16
 * upgrade to the latest release of protobuf

v1.0.43
----------
 * Update metadata for v8.10.14

v1.0.42
----------
 * Update for metadata changes in v8.10.13
 * fix yoda expressions
 * fix slice operations
 * fix regex escaping
 * fix make calls
 * fix error strings

v1.0.41
----------
 * update metadata for v8.10.12

v1.0.40
----------
 * add unit test for valid/possible US/CA number, include commit in netlify version, lastest metadata
 * update readme to add svn dependency

v1.0.39
----------
 * add dist to gitignore
 * tweak goreleaser

v1.0.38
----------
 * update travis env to always enable modules

v1.0.37
----------
 * plug in goreleaser and add it to travis

v1.0.36
----------
 * Update for upstream metadata v8.10.7

v1.0.35
----------
 * update metadata for v8.10.4 release
 * update AR test number to valid AR fixed line

v1.0.34
----------
 * update travis file

v1.0.33
----------
 * remove goreleaser since we no longer use docker for test deploys
 * latest google metadata

v1.0.32
----------
 * add /functions to gitignore
 * update to latest google metadata

v1.0.31
----------
 * update to latest metadata v8.10.1, test case changes validated against google lib
 * add link in readme to test function

v1.0.30
----------
 * fix FormatByPattern with user defined pattern. Fixes: #16

v1.0.29
----------
 * update metadata v8.9.16 (test diff validated against python lib)

v1.0.28
----------
 * update metadata to v8.9.14, fix go.mod dependency

v1.0.27
----------
 * update to metadata v8.9.13, remove must dependency

v1.0.26
----------
 * Fix cache strict look up bug and unify cache management, thanks @eugene-gurevich

v1.0.25
----------
 * save possible lengths to metadata, change implementation to use, add IS_POSSIBLE_LOCAL_ONLY and INVALID_LENGTH as possible return values to IsPossibleNumberWithReason
 * update metadata to version v8.9.12

v1.0.24
----------
 * update to metadata for v8.9.10

v1.0.23
----------
 * add GetSupportedCallingCodes
 * return sets as map[int]bool instead of map[int]struct{}

v1.0.22
----------
* add GetCarrierForNumber and GetGeocodingForNumber

v1.0.21
----------
 * Update for libphonenumber v8.9.8

v1.0.20
----------
 * updated metadata for v8.9.7

v1.0.19
----------
 * update metadata for v8.9.6

v1.0.18
----------
 * update metadata for v8.9.5

v1.0.17
----------
 * Fix maybe strip extension, thanks @vlastv

