import os
import json
import sys
import requests
import fdk
import ssl
import ujson

from urllib import request
from google.oauth2 import service_account
from google.cloud import vision
from google.cloud.vision import types


def inject_client(client):

    def handle(context, data=None, loop=None):
        try:
            ctx = None
            data = ujson.loads(data)
            image_url = data.get("media_url")
            if "https" in image_url:
                ctx = ssl.create_default_context()
                ctx.check_hostname = False
                ctx.verify_mode = ssl.CERT_NONE

            print("Image URL: ", image_url, file=sys.stderr, flush=True)
            user = data.get("user")
            tweet_id = data.get("tweet_id")

            try:
                url_response = request.urlopen(image_url, context=ctx)
                content = url_response.read()
            except Exception as ex:

                print("SHIT 1", str(ex), flush=True, file=sys.stderr)
                tweet_fail = data.get("tweet_fail")
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
                    print(landmark, file=sys.stderr, flush=True)
                    tweet_success = data.get("tweet_success")
                    requests.post(tweet_success, json={
                        "user": user,
                        "tweet_id": tweet_id,
                        "landmark": landmark,
                    })
            else:
                print("SHIT 2", flush=True, file=sys.stderr)
                tweet_fail = data.get("tweet_fail")
                requests.post(tweet_fail, json={
                    "user": user,
                    "tweet_id": tweet_id,
                })
        except Exception as ex:
            print("SHIT 3", flush=True, file=sys.stderr)
            sys.stderr.write("STDIN data: {}\n".format(
                json.dumps(data, sort_keys=False, indent=4)))
            sys.stderr.write("ENV: {}\n".format(
                json.dumps(dict(os.environ), sort_keys=False, indent=4)))
            sys.stderr.write("----------------------------\n{}"
                             "\n----------------------------"
                             .format(str(ex)))
            raise ex

    return handle


def setup_client():
    g_type = os.environ.get("TYPE", os.environ.get("type"))
    g_project_id = os.environ.get("PROJECT_ID", os.environ.get("project_id"))
    g_private_key_id = os.environ.get("PRIVATE_KEY_ID", os.environ.get("private_key_id"))
    g_private_key = os.environ.get("PRIVATE_KEY", os.environ.get("private_key"))
    g_client_email = os.environ.get("CLIENT_EMAIL", os.environ.get("client_email"))
    g_client_id = os.environ.get("CLIENT_ID", os.environ.get("client_id"))
    g_auth_uri = os.environ.get("AUTH_URI", os.environ.get("auth_uri"))
    g_token_uri = os.environ.get("TOKEN_URI", os.environ.get("token_uri"))
    g_auth_provider_x509_cert_url = os.environ.get(
        "AUTH_PROVIDER_X509_CERT_URL", os.environ.get("auth_provider_x509_cert_url"))
    g_client_x509_cert_url = os.environ.get(
        "CLIENT_X509_CERT_URL", os.environ.get("client_x509_cert_url"))
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
    return vision.ImageAnnotatorClient(
        credentials=credentials,
        scopes=['https://www.googleapis.com/auth/cloud-platform', ])


if __name__ == "__main__":
    fdk.handle(inject_client(setup_client()))
