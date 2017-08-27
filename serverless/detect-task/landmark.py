import os
import json
import sys
import requests

from urllib import request
from google.oauth2 import service_account
from google.cloud import vision
from google.cloud.vision import types

DISCOVERY_URL = 'https://{api}.googleapis.com/$discovery/rest?version={apiVersion}'

# {"media_url": "https://pbs.twimg.com/media/DIL4q-3WsAAEXia.jpg:large", "user": "@denis_makogon", "tweet_id": "901556605967380480"}

# TODO(denismakogon): better error handling
if __name__ == "__main__":
    if not os.isatty(sys.stdin.fileno()):
        try:
            g_type = os.environ.get("type")
            g_project_id = os.environ.get("project_id")
            g_private_key_id = os.environ.get("private_key_id")
            g_private_key = os.environ.get("private_key")
            g_client_email = os.environ.get("client_email")
            g_client_id = os.environ.get("client_id")
            g_auth_uri = os.environ.get("auth_uri")
            g_token_uri = os.environ.get("token_uri")
            g_auth_provider_x509_cert_url = os.environ.get("auth_provider_x509_cert_url")
            g_client_x509_cert_url = os.environ.get("client_x509_cert_url")

            if not all([g_type, g_project_id, g_private_key_id, g_private_key,
                        g_client_email, g_auth_uri, g_token_uri,
                        g_auth_provider_x509_cert_url, g_client_x509_cert_url]):
                raise Exception("One or more GCloud auth attributes empty.")

            g_private_key = g_private_key.replace("\\n", "\n")
            gcloup_map = {
                "type": g_type,
                "project_id": g_project_id,
                "private_key_id": g_private_key_id,
                "private_key": g_private_key,
                "client_email": g_client_email,
                "client_id": g_client_id,
                "auth_uri": g_auth_uri,
                "token_uri": g_token_uri,
                "auth_provider_x509_cert_url": g_auth_provider_x509_cert_url,
                "client_x509_cert_url": g_client_x509_cert_url,
            }

            credentials = service_account.Credentials.from_service_account_info(
                gcloup_map, scopes=['https://www.googleapis.com/auth/cloud-platform', ])
            client = vision.ImageAnnotatorClient(
                credentials=credentials,
                scopes=['https://www.googleapis.com/auth/cloud-platform', ])

            obj = json.loads(sys.stdin.read())
            image_url = obj.get("media_url")
            if not image_url:
                # TODO(denismakogon): tweet back with bad image URL
                sys.stderr.write("Empty media URL")
                raise Exception("Empty media URL")
            # need to download image, remote image URI is not stable
            user = obj.get("user")
            tweet_id = obj.get("tweet_id")
            content = None
            try:
                filename, _ = request.urlretrieve(image_url)
                with open(filename, 'rb') as image_file:
                    content = image_file.read()
            except Exception as ex:
                # TODO(denismakogon): tweet with bad image URL
                tweet_fail = obj.get("tweet_fail")
                requests.post(tweet_fail, json={
                    "user": user,
                    "tweet_id": tweet_id,
                    "bad_image_source":  True,
                })
                raise ex
            image = types.Image(content=content)
            response = client.landmark_detection(image=image)
            landmarks = response.landmark_annotations

            if len(landmarks) > 0:
                possible_landmarks = set(
                    [landmark.description for landmark in landmarks])
                for landmark in possible_landmarks:
                    tweet_success = obj.get("tweet_success")
                    requests.post(tweet_success, json={
                        "user": user,
                        "tweet_id": tweet_id,
                        "landmark": landmark,
                    })
            else:
                tweet_fail = obj.get("tweet_fail")
                requests.post(tweet_fail, json={
                    "user": user,
                    "tweet_id": tweet_id,
                })
        except Exception as ex:
            sys.stderr.write(str(ex))
