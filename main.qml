import QtQuick 2.2
// import QtQuick.Controls 1.1
// import GoExtensions 1.0

Item {
    id: root
    width: 1024
    height: 768

    property list<View> views


    // TODO: Implement exclusivegroup for views so that I can change visibility more carefully
    MainMenu {
        id: mainMenuView
        name: "mainMenuView"
        visible: true
    }

    MapView {
        id: mapView
        name: "gameView"
        visible: false
    }
}