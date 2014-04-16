import QtQuick 2.1
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0
// import GoExtensions 1.0


Dialog {
    id: connectDialog
    width: 300
    height: 80

    ColumnLayout {
        id: mainLayout
        anchors.fill: parent
        anchors.margins: margin

        RowLayout {
            id: rowLayout
            Layout.fillWidth: true

            TextField {
                id: addressTextField
                placeholderText: "Server address..."
                text: "127.0.0.1:2222"
                Layout.fillWidth: true
            }
        }

        RowLayout {
            id: buttonsRow
            Layout.fillWidth: true
            EpicButton {
                text: "Connect"
                width: 100
                height: 30
                onClicked: {
                    game.address = addressTextField.text
                    console.log("Address: ", addressTextField.text)
                    game.state = 1
                    if (game.state == 1) {
                        mainMenuView.visible = false
                        mapView.visible = true
                    }
                }
            }
            EpicButton {
                text: "Cancel"
                width: 100
                height: 30
                onClicked: {
                    connectDialog.visible = false
                }
            }
        }

    }
}