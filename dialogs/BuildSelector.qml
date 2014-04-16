import QtQuick 2.1

Rectangle {
    id: buildSelector

    radius: 10

    property string selected_archetype: "blank32"

    border.width: 1
    border.color: "gray"

    Component {
        id: buildDelegate
        Item {
            width: grid.cellWidth; height: grid.cellHeight
                
            Image {
                id: csImage
                source: {
                    return "../art/" + buildModel.imageName(index) + ".png"
                }
            }

            Text {
                anchors.right: parent.right; anchors.bottom: parent.bottom
                style: Text.Outline; styleColor: "#AAAAAA"
                text: buildModel.get(index).quantity

                Component.onCompleted: {
                    if (text == "0") {
                        text = ""
                    }
                }
            }

            MouseArea {
                id: mouseArea
                acceptedButtons: Qt.LeftButton
                width: parent.width
                height: parent.height
                hoverEnabled: true         //this line will enable mouseArea.containsMouse

                onClicked: {
                    buildSelector.selected_archetype = buildModel.get(index).archetype.simpleName
                    game.selectedArchetype = buildSelector.selected_archetype
                    console.log("Index " + index + " clicked to build.  Archetype set to " + game.selectedArchetype)
                    buildSelector.visible = false
                }
            }
        }
    }

    GridView {
        id: grid
        interactive: false
        anchors.fill: parent
        cellWidth: 34; cellHeight: 34

        model: buildModel.size
        delegate: buildDelegate
    }
}