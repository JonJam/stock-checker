# Deployment
This is the setup to deploy to Digital Ocean.

## App Platform
### Steps
1. If app not created, to create run `doctl apps create --spec app.yml`
2. If app exists, to update run `doctl apps update [INSERT APP ID] --spec app.yml`
3. Add Twilio environment variables with secrets by following [this](https://www.digitalocean.com/docs/app-platform/how-to/use-environment-variables/#define-build-time-environment-variables). When specifying values, do not use `"`.
```
SC_TWILIO_ENABLED
SC_TWILIO_ACCOUNTSID
SC_TWILIO_AUTHTOKEN
SC_TWILIO_NUMBERTO
SC_TWILIO_NUMBERFROM
```