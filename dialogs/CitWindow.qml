import QtQuick 2.1
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0


Item {
    id: rawDialog
    property int margin: 11

    property string title: "None"

	property int content_x: 10
	property int content_width: width - (content_x * 2)
	property int content_height: height - content_y - 10
    property int content_y: headerBar.height + 10

    z: 0

    Rectangle {
        width: parent.width
        height: parent.height
        radius: 10
        color: "#cccccc"

	    MouseArea {
	    	// Exists to stop mouse events passing through
	    	acceptedButtons: Qt.LeftButton | Qt.RightButton
	    	anchors.fill: parent
	    	onClicked: {
	    		console.log("Z was: " + rawDialog.z)
		    	rawDialog.z = game.newWindowZ(z) // For now, just returns one higher z, and sets its record to highest
		    	console.log("Game z set to: " + rawDialog.z)
		    }
	    }

        Rectangle {
        	id: headerBar
	        z: parent.z + 1
	        width: parent.width
	        height: 15
	        color: "#eeeeee"
	        radius: 10
	        Text {
	        	anchors.centerIn: parent
		        text: title
		    }

        	MouseArea {
		       id: mouseArea
		       acceptedButtons: Qt.LeftButton | Qt.RightButton
		       anchors.fill: parent
		       width: parent.width
		       height: parent.height

		       drag.target: rawDialog
		       drag.axis: Drag.XandYAxis
		       drag.minimumX: 0
		       drag.minimumY: 0
		       drag.maximumX: rawDialog.parent.width - rawDialog.width
		       drag.maximumY: rawDialog.parent.height - rawDialog.height
			}
	    }
    }

}