import QtQuick 2.1
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0
// import GoExtensions 1.0


CitWindow {
    id: craftingWindow
    width: 700
    height: 350

    title: "Crafting"

    Rectangle {
        id: craftingRect
        clip: true
        x: content_x
        y: content_y
        width: content_width
        height: content_height

        property int selected: 0
        property var craftElement

        onSelectedChanged: {

        }

        Item {
            anchors.left: parent.left
            anchors.top: parent.top
            width: parent.width/2
            height: parent.height
            Component {
                id: craftDelegate
                Item {
                    width: parent.width; height: 40

                    Rectangle {
                        radius: 5
                        color: "gray"
                        border.width: 1
                        border.color: "#222222"
                        opacity: 0.5
                        anchors.fill: parent
                    }

                    Row {
                        Item {
                            width: 45
                            height: parent.parent.height
                            Image {
                                anchors.verticalCenter: parent.verticalCenter
                                source: "../art/" + craftModel.get(index).archetype.simpleName + ".png" }
                            }
                        Text { text: craftModel.get(index).archetype.name }
                    }
                    MouseArea {
                        id: mouseArea
                        acceptedButtons: Qt.LeftButton | Qt.RightButton
                        width: parent.width
                        height: parent.height
                        hoverEnabled: true         //this line will enable mouseArea.containsMouse

                        onClicked: {
                            craftingRect.selected = index
                            craftingRect.craftElement = craftModel.get(index)
                        }
                    }
                }
            }

            ListView {
                anchors.fill: parent
                model: craftModel.len
                delegate: craftDelegate
                highlight: Rectangle { color: "lightsteelblue"; radius: 5 }
                focus: true
            }
        }

        Item {
            id: selectedView
            anchors.right: parent.right
            anchors.top: parent.top
            width: parent.width/2
            height: parent.height
            Rectangle {
                id: headerRow
                radius: 5
                border.width: 1
                border.color: "#222222"
                // Header info for item
                width: parent.width
                height: 48
                Image {
                    anchors.verticalCenter: parent.verticalCenter
                    source: "../art/" + craftingRect.craftElement.archetype.simpleName + ".png"
                }
                Text {
                    anchors.horizontalCenter: parent.horizontalCenter
                    text: craftingRect.craftElement.archetype.name 
                }
            }

            Column {
                anchors.top: headerRow.bottom

                Text {
                    id: reqHeader
                    text: "Requirements:"
                }

                Component {
                    id: craftReqsDelegate
                    Item {
                        width: parent.width; height: 40


                        Row {
                            Item {
                                width: 45
                                height: parent.parent.height
                                Image {
                                    anchors.verticalCenter: parent.verticalCenter
                                    source: "../art/" + craftingRect.craftElement.get(index).archetype.simpleName + ".png" }
                                }
                            Text { text: craftingRect.craftElement.get(index).quantity + " x " + craftingRect.craftElement.get(index).archetype.name }
                        }

                    }
                }

                ListView {
                    anchors.top: reqHeader.bottom
                    
                    // anchors.fill: parent
                    model: craftingRect.craftElement.len
                    delegate: craftReqsDelegate
                    highlight: Rectangle { color: "lightsteelblue"; radius: 5 }
                    focus: true
                }
            }

            EpicButton {
                text: "craft!"
                width: 80
                height: 30
                shadow: false
                anchors.bottom: parent.bottom
                anchors.horizontalCenter: parent.horizontalCenter
                onClicked: {
                    var be = craftModel.get(craftingRect.selected)
                    game.createJob("craft", false, {"product_name": be.archetype.simpleName})
                }
            }
        }
    }
}