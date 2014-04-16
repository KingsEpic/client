import QtQuick 2.2
import QtQuick.Controls 1.1
import "dialogs"

View {
    width: root.width
    height: root.height


    Image {
        anchors.fill: parent
        source: "art/title.png"
        fillMode: Image.PreserveAspectCrop
        smooth: true
    }

    Item {
        x: 10
        y: parent.height - 10
        Column {
            id: menu_buttons
            width: childrenRect.width
            anchors.bottom: parent.bottom
            anchors.left: parent.left
            spacing: 8
            
            EpicButton { text: "Connect"; onClicked: {connect_dialog.visible = (connect_dialog.visible? false:true)} width: 100; height: 40;}
            EpicButton { text: "About"; onClicked: {about_dialog.visible = (about_dialog.visible? false:true)} width: 100; height: 40;}
        }
    }

    About {
        id: about_dialog
        anchors.centerIn: parent
        visible: false
    }

    Connect {
        id: connect_dialog
        anchors.centerIn: parent
        visible: false
    }

}