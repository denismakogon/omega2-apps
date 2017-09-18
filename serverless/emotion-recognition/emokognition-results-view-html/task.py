import jinja2
import requests
import os
import sys


if __name__ == "__main__":
    loader = jinja2.FileSystemLoader('./index.html')
    env = jinja2.Environment(loader=loader)
    template = env.get_template('')

    fn_app = os.environ.get("FN_APP_NAME")
    recorder = "{0}/r/{1}/results".format(os.environ.get("FN_API_URL"), fn_app)
    try:
        resp = requests.get(recorder)
        resp.raise_for_status()
        data = resp.json()
        main_emotions = data['main']
        alt_emotions = data['alt']
        main, alt = [], []
        total = sum(list(main_emotions.values()))
        for emotion, count in main_emotions.items():
            main.append(dict(emotion=emotion,
                             stat=float("{:.2f}".format(float(count / total) * 100)),
                             times=count))
        for emotion, count in alt_emotions.items():
            alt.append(dict(emotion=emotion,
                            stat=float("{:.2f}".format(float(count / total) * 100)),
                            times=count))
        context = {
            "main_emotions": main,
            "alt_emotions": alt,
            "total": total
        }
        print(template.render(context))
    except Exception as ex:
        sys.stderr.write(str(ex))
