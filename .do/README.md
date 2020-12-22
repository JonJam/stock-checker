# Deployment
This is the setup to deploy to Digital Ocean.

## App Platform
### Steps
1. If app not created, run `doctl apps create --spec app.yml`
2. If app exists, run `doctl apps update [INSERT APP ID] --spec app.yml`