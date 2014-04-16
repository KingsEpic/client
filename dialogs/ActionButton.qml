import QtQuick 2.1
import QtQuick.Controls 1.1

Item {
	id: actionButton
	width: 64
	height: 64

	property bool checked: false
	property bool allowNone: false
	property ExclusiveGroup exclusiveGroup: null

	property string iconSource

	onExclusiveGroupChanged: {
		if (exclusiveGroup) 
			exclusiveGroup.bindCheckable(actionButton)
	}

	onCheckedChanged: {
		// console.log("Inside: check changed")
		if (checked) {
			actionIcon.scale = 0.8
		} else {
			actionIcon.scale = 1.0
		}
	}

	Image {
    	id: actionIcon
    	width: parent.width; height: parent.height;
    	fillMode: Image.PreserveAspectFit
    	smooth: true
        source: iconSource
    }

	MouseArea {
       id: mouseArea
       width: parent.width
       height: parent.height

       onClicked: {
       	if (!allowNone) {
       		checked = true
       	} else {
	       	if (checked) {
	       		checked = false
	       	} else {
		        checked = true
		    }
		}
      }
	}
}