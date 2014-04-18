import QtQuick 2.1
import QtQuick.Controls 1.1

Item {
	id: epicButton

	property string text
	signal clicked()
  property bool shadow: true

  onShadowChanged: {
    shadowRect.visible = shadow
  }


  Rectangle {
    id: shadowRect
    anchors.centerIn: parent
    width: parent.width + 6
    height: parent.height + 6
    radius: 5
    opacity: 0.4
    color: "#222222"
  }

	Image {
    	id: actionIcon
    	width: parent.width; height: parent.height;
    	smooth: true
        source: "../art/button.png"
  }

  Text {
  	anchors.centerIn: parent
    color: "white"
    font.bold: true
    style: Text.Raised;
    font.capitalization: Font.AllUppercase
  	text: parent.text
  }

	MouseArea {
       id: mouseArea
       width: parent.width
       height: parent.height
       onClicked: {
        epicButton.clicked()
       }

       onPressed: {
        parent.scale = 0.9
       }

       onReleased: {
        parent.scale = 1.0
       }
	}
}