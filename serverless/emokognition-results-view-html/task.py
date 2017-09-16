import jinja2
import requests
import os
import sys


if __name__ == "__main__":
    loader = jinja2.FileSystemLoader('./index.html')
    env = jinja2.Environment(loader=loader)
    template = env.get_template('')

    recorder = "{0}/r/emokognition/results".format(os.environ.get("FN_API_URL"))
    try:
        resp = requests.get(recorder)
        resp.raise_for_status()
        data = resp.json()
        main_emotions = data['main']
        alt_emotions = data['alt']
        main, alt = [], []
        total = sum(list(main_emotions.values()))
        for emotion, count in main_emotions.items():
            main.append(dict(emotion=emotion, stat=float(count / total) * 100))
        for emotion, count in alt_emotions.items():
            alt.append(dict(emotion=emotion, stat=float(count / total) * 100))
        context = {
            "main_emotions": main,
            "alt_emotions": alt,
        }
        print(template.render(context))
    except Exception as ex:
        sys.stderr.write(str(ex))
