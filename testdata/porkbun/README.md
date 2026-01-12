# Solver testdata directory

## Running Tests

To run the conformance tests, you need to provide valid Porkbun API credentials:

### Understanding the test files:

- **config.json** - Contains references to where credentials are stored. DO NOT modify this file - it tells the webhook to look for keys named "apiKey" and "secretApiKey" in the Secret named "porkbun-credentials".

- **porkbun-credentials.yaml** - Contains the actual base64-encoded credentials. THIS is the file you need to update with your API keys.

### Setup steps:

1. **Enable API Access for your domain in Porkbun:**
   - Log in to Porkbun and go to your domain management page
   - Find the domain you want to test with (e.g., `example.com`)
   - Look for "API Access" setting and make sure it's **ENABLED**
   - **Important:** Without enabling API Access for the specific domain, Porkbun will return ERROR status for all API calls, even with valid credentials

2. Get your Porkbun API credentials from https://porkbun.com/account/api
3. Base64 encode your credentials:
   ```bash
   echo -n "pk1_your_api_key" | base64
   echo -n "sk1_your_secret_key" | base64
   ```
4. Edit `porkbun-credentials.yaml` and replace the placeholder values:
   - Replace `<your-api-key-base64-encoded>` with the output from the first command
   - Replace `<your-secret-api-key-base64-encoded>` with the output from the second command

   **Important:** Put the base64-encoded values in `porkbun-credentials.yaml`, NOT in `config.json`!

5. Set the TEST_ZONE_NAME environment variable to a domain you own in Porkbun
6. Run the tests:
   ```bash
   TEST_ZONE_NAME=yourdomain.com. make test
   ```

Note: The test will create and delete TXT records in your actual DNS zone.
