from model import *


def main():
    filename = 'models/trained_model.sav'

    # train model
    train_data = get_train_data()
    x_train, x_test, y_train, y_test = train_test_split(0.8, train_data)
    model = KNeighborsClassifier()
    train_model(x_train, y_train, x_test, y_test, model, filename)


if __name__ == '__main__':
    main()
