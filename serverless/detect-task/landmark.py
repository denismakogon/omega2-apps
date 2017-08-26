from urllib import request
import os
import json
import sys


from googleapiclient import discovery
from googleapiclient import errors
from oauth2client.client import GoogleCredentials


DISCOVERY_URL = 'https://{api}.googleapis.com/$discovery/rest?version={apiVersion}'


if __name__ == "__main__":
    if not os.isatty(sys.stdin.fileno()):
        try:
            # TODO(denismakogon):Let func consumer gclouth auth through config instead of stdin
            obj = json.loads(sys.stdin.read())
            gcloup_map = obj.get('gcloud')
            json.dumps(gcloup_map)
            credentials = GoogleCredentials.from_json(json.dumps(gcloup_map))
            # credentials = GoogleCredentials.get_application_default()
            service = discovery.build(
                'vision', 'v1', credentials=credentials,
                discoveryServiceUrl=DISCOVERY_URL)

            image_url = obj.get("media_url")
            response = request.urlopen(image_url)
            data = response.read()
            data.decode('utf-8')

            batch_request = [{
                'image': {
                    'content': data,
                },
                'features': [{
                    'type': 'LANDMARK_DETECTION',
                    'maxResults': 4,
                }]
            }]
            request = service.images().annotate(body={'requests': batch_request})
            responses = request.execute(num_retries=1)
            if 'responses' not in responses:
                pass
                # do something here
            if 'error' in responses:
                if 'message' in responses['error']:
                    sys.stderr.write(response['error']['message'])
            text_response = {}

        except errors.HttpError as ex:
            sys.stderr.write("Http Error %s" % str(ex))
        except Exception as ex:
            sys.stderr.write("Error: %s" % str(ex))

    # TODO(denismakogon): landmark detection should happen here
    # TODO(denismakogon): use config-based auth, but not request data-based
