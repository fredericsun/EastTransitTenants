import pandas as pd
import numpy as np
import glob
import sys
import json

from sklearn.linear_model import LogisticRegression
from sklearn.neighbors import KNeighborsClassifier
from sklearn.svm import SVC


def get_train_data():
    all_json = glob.glob("data/ticket_*.json")
    all_data = [pd.read_json(filename) for filename in all_json]
    train_data = pd.concat(all_data)
    for filename in all_json:
        df = pd.read_json(filename)
    train_data = train_data.sample(frac=1).reset_index(drop=True)
    # train_data = train_data[['TotalPercent', 'CountPercent', 'TotalService', 'Load', 'BottleNeck']]
    return train_data


def train_test_split(ratio, train_data):
    split = int(ratio * len(train_data))
    x_train = train_data.iloc[:split, :-1].values
    x_test = train_data.iloc[split:, :-1].values
    y_train = train_data.iloc[:split, -1].values
    y_test = train_data.iloc[split:, -1].values
    return x_train, x_test, y_train, y_test


def train_model(x_train, y_train, x_test, y_test, model):
    model.fit(x_train, y_train)
    acc = model.score(x_test, y_test)
    return model, acc


def main():
    train_data = get_train_data()
    x_train, x_test, y_train, y_test = train_test_split(0.8, train_data)
    model = KNeighborsClassifier()
    fitted_model, acc = train_model(x_train, y_train, x_test, y_test, model)

    input_path = pd.read_json("data/path.json")
    data = input_path.iloc[:, 1:].values
    names = input_path.iloc[:, 0].values
    result = [bool(ele) for ele in fitted_model.predict(data)]
    output = {names[i]: result[i] for i in range(len(names))}
    with open("data/output.json", "w") as outfile:
        json.dump(output, outfile)


if __name__ == '__main__':
    main()
