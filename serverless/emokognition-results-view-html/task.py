import jinja2
import requests
import os

loader = jinja2.FileSystemLoader('./index.html')
env = jinja2.Environment(loader=loader)
template = env.get_template('')


if __name__ == "__main__":
    api_url = os.environ.get("FN_API_URL")
    recorder = "{}/r/emokognition/recorder".format(api_url)
    data = requests.get(recorder)
    items, alt_items = [], []
    main_emotions = data['main_emotion']
    alt_emotions = data['alt_emotion']
    total = sum(list(main_emotions.values()))
    for emotion, count in main_emotions.items():
        items.append(dict(emotion=emotion, stat=float(count / total) * 100))
    for emotion, count in alt_emotions.items():
        alt_items.append(dict(emotion=emotion, stat=float(count / total) * 100))

    print(template.render(items=items, alt_items=alt_items))
