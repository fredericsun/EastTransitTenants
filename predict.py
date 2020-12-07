from model import *
import joblib
import pandas as pd


def main():
    filename = 'models/trained_model.sav'

    input_path = pd.read_json("path/path.json")
    input_path = input_path.groupby(['ServiceName']).mean(
    )[['SelfPercent', 'ServiceDur', 'TotalPercent', 'CountPercent', 'TotalDur', 'TotalService', 'TotalSpan', 'Load']]

    data = input_path.iloc[:, :].values
    names = input_path.index.values

    loaded_model = joblib.load(filename)
    result = [bool(ele) for ele in loaded_model.predict(data)]
    output = {names[i]: result[i] for i in range(len(names))}

    with open("path/output.json", "w") as outfile:
        json.dump(output, outfile)


if __name__ == '__main__':
    main()
