import requests
import base64
import time
import json
import yaml
import pprint
import os
languages = {
    # "python":{
    #     "frameworks":["run","unittest"],
    #     "timeout":20000
    # },
    "golang":{
        "frameworks":["run"],
         "timeout":20000
    }
}

bs4_encode = lambda x:base64.encodebytes(x.encode()).decode("utf-8")

file_exists =  lambda f:os.path.exists(f)

def key_exists(d,k) -> bool:
    try:
        d[k]
        return True
    except KeyError:
        return False

for lang, v in languages.items():
    print("Language:%s"%lang)
    if not file_exists("{lang}.yml".format(lang=lang)):
        continue
    language_examples = None 
    with open("{lang}.yml".format(lang=lang),"r") as outfile:
        language_examples = yaml.load(outfile,Loader=yaml.FullLoader)
    for f in v["frameworks"]:
        if not key_exists(language_examples,f):
            continue
        print("Framework:%s"%f)
        for k in language_examples[f].keys():
            print("Test Name:%s"%k)
            test = language_examples[f][k]
            payload = {
                "language":lang,
                "code":bs4_encode(test["code"]),
            }
            if key_exists(test,"setup"):
                payload["setup"]=bs4_encode(test["setup"])
            if key_exists(test,"fixture"):
                payload["fixture"]=bs4_encode(test["fixture"])
            if key_exists(v,"timeout"):
                payload["timeout"]=v["timeout"]
            start = time.time()
            r = requests.post("http://127.0.0.1:8000/test",json=payload)
            print("Time taken: ",time.time()-start)
            if r.status_code == 200:
                try:
                    pprint.pprint(r.json())
                except json.decoder.JSONDecodeError:
                    print(r.text)
            else:
                print(r.text)