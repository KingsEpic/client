import QtQuick 2.1

Item {
    id: inventory

    property int selected_index: 0

    Rectangle {
        anchors.fill: parent
        opacity: 0.5
    }

    Component {
        id: invDelegate
        Item {

            width: grid.cellWidth; height: grid.cellHeight
                
            Image {
                id: csImage

                source: {
                    return "../art/" + invModel.get(index).imageName() + ".png"
                }
            }

            Text {
                anchors.right: parent.right; anchors.bottom: parent.bottom
                style: Text.Outline; styleColor: "#AAAAAA"
                text: invModel.get(index).quantity()

                Component.onCompleted: {
                    if (text == "0") {
                        text = ""
                    }
                }

            }

            MouseArea {
                id: mouseArea
                acceptedButtons: Qt.LeftButton | Qt.RightButton
                width: parent.width
                height: parent.height
                hoverEnabled: true         //this line will enable mouseArea.containsMouse

                onClicked: {
                    console.log("Clicked inv item: " + invModel.get(index).index + " the id: " + invModel.get(index).id +  " entity: " + invModel.get(index).entity.simpleName)
                    if (inventory.selected) {
                        // request tile swap:
                        inventory.selected = false
                        itemRect.visible = false
                        // console.log("Requested swap item in " + inventory.selected_index + " to this index at " + index)
                        // invModel.swap(inventory.selected_index, index)
                    } else {
                        inventory.selected = true
                        itemRect.visible = true
                        inventory.selected_index = index
                        // console.log("Selected this index " + index)
                    }
                }
            }
        }
    }

    GridView {
        id: grid
        interactive: false
        anchors.fill: parent
        cellWidth: 34; cellHeight: 34

        model: invModel.size
        delegate: invDelegate
    }
}