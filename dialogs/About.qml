import QtQuick 2.0


Dialog {
    id: aboutComponent
    width: 300
    height: 300

    MouseArea {
        anchors.fill: parent
        onClicked: { aboutComponent.visible = false} 
    }

    Item{
        width: parent.width - 50
        height: parent.height - 50
        anchors.centerIn: parent
        Text {
        	anchors.fill: parent
            color: "white"
            font.bold: true
            style: Text.Raised;
            font.capitalization: Font.AllUppercase
            wrapMode: Text.WordWrap
        	text: "King's Epic.
Copyright by Mark Saward.

Please see the provided LICENSE file for details."
        }
    }
}