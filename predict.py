from model import *
import joblib
import pandas as pd


def main():
    filename = 'models/trained_model.sav'

    # predict
    input_path = pd.read_json("data/path.json")
    data = input_path.iloc[:, 1:].values
    names = input_path.iloc[:, 0].values

    loaded_model = joblib.load(filename)
    result = [bool(ele) for ele in loaded_model.predict(data)]
    output = {names[i]: result[i] for i in range(len(names))}
    with open("data/output.json", "w") as outfile:
        json.dump(output, outfile)


if __name__ == '__main__':
    main()
