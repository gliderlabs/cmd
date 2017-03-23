// Package dynamodb implements a dynamodb store backend with migration support.
//
// Deploying Hard Migrations
//
// Maintenance mode MUST be active in order for CI deployment to succeed.
//
//  1. Activate maintenance mode for release channel within your PR.
//  2. Merge PR and wait for deployment to succeed.
//  3. Manually apply migrations.
//  4. Create and merge a second PR disabling maintenance.
//
// Applying Hard Migrations
//
// Hard migration can be applied by executing cmd with the "-migrate" flag and
// maintenance mode active.
//
//  # Local
//  MAINTENANCE_ACTIVE=true ./build/darwin_amd64/cmd -migrate dev/dev.toml
//
//  # Production against alpha channel (assuming config spec activated maintenance)
//  kubectl exec -it $PODNAME -- /usr/local/bin/cmd -d /config/config.toml -migrate
//
package dynamodb
