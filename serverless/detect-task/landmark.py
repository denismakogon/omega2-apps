import os
import json
import sys
import requests

from urllib import request
from google.oauth2 import service_account
from google.cloud import vision
from google.cloud.vision import types


if __name__ == "__main__":
    if not os.isatty(sys.stdin.fileno()):
        try:
            sys.stderr.write("ENV: {}\n".format(
                json.dumps(dict(os.environ), sort_keys=False, indent=4)))
            g_type = os.environ.get("TYPE")
            g_project_id = os.environ.get("PROJECT_ID")
            g_private_key_id = os.environ.get("PRIVATE_KEY_ID")
            g_private_key = os.environ.get("PRIVATE_KEY")
            g_client_email = os.environ.get("CLIENT_EMAIL")
            g_client_id = os.environ.get("CLIENT_ID")
            g_auth_uri = os.environ.get("AUTH_URI")
            g_token_uri = os.environ.get("TOKEN_URI")
            g_auth_provider_x509_cert_url = os.environ.get("AUTH_PROVIDER_X509_CERT_URL")
            g_client_x509_cert_url = os.environ.get("CLIENT_X509_CERT_URL")

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
            sys.stderr.write("GCloud map: {}\n".format(
                json.dumps(gcloup_map, sort_keys=False, indent=4)))
            credentials = service_account.Credentials.from_service_account_info(
                gcloup_map, scopes=['https://www.googleapis.com/auth/cloud-platform', ])
            client = vision.ImageAnnotatorClient(
                credentials=credentials,
                scopes=['https://www.googleapis.com/auth/cloud-platform', ])

            obj = json.loads(sys.stdin.read())
            image_url = obj.get("media_url")
            user = obj.get("user")
            tweet_id = obj.get("tweet_id")
            content = None
            try:
                filename, _ = request.urlretrieve(image_url)
                with open(filename, 'rb') as image_file:
                    content = image_file.read()
            except Exception as ex:
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
                sys.stderr.write("Possible landmarks: {}\n"
                                 .format(possible_landmarks))
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
            sys.exit(0)
