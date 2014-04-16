import QtQuick 2.1
import QtQuick.Controls 1.1

Item {
	id: actionBar
	ExclusiveGroup { id: actionButtons }

    x: width - 72
    y: 8

    BuildSelector {
    	id: buildSelector
    	width: 200
    	height: 200
    	x: -200
    	visible: false
    }

	Column {
        spacing: 2

        Row {
        	ExclusiveGroup { id: delegationMode }

        	ActionButton {
		        iconSource: "../art/icons/no_delegate.png"
		        checked: true
		        width: 32; height:32
		        exclusiveGroup: delegationMode
		        onCheckedChanged: {
		        	if (checked) {
		        		game.delegated = false
		        	}
		        }
		    }

		    ActionButton {
		        iconSource: "../art/icons/delegate.png"
		        width: 32; height:32
		        exclusiveGroup: delegationMode
		        onCheckedChanged: {
		        	if (checked) {
		        		game.delegated = true
		        	}
		        }
		    }
        }

		ActionButton {
	        iconSource: "../art/icons/walk.png"
	        checked: true
	        exclusiveGroup: actionButtons
	        onCheckedChanged: {
	        	if (checked) {
	        		actionBar.parent.right_action = "walk"
	        	}
	        }
	    }

	    ActionButton {
	        iconSource: "../art/icons/harvest.png"
	        exclusiveGroup: actionButtons
	        onCheckedChanged: {
	        	if (checked) {
	        		actionBar.parent.right_action = "harvest"
	        	}
	        }
	    }

	   	Row {
	   		id: buildRow
	   		x: -32
	   		Item {
	   			width: 32; height: 32
	   			anchors.verticalCenter: parent.verticalCenter

	   			Rectangle {
	   				anchors.fill: parent
		   			opacity: 0.5
		   			
		   		}

		       	Image {
			    	id: buildIcon
			    	width: parent.width; height: parent.height;
			        source: "../art/" + buildSelector.selected_archetype + ".png"
			    }

	   			MouseArea {
			        anchors.fill: parent

			        onClicked: {
			        	buildSelector.visible = (buildSelector.visible ? false:true)
			        }
			    }
		    }



		    ActionButton {
		        iconSource: "../art/icons/build.png"
		        exclusiveGroup: actionButtons
		        onCheckedChanged: {
		        	if (checked) {
		        		actionBar.parent.right_action = "build"
		        	}
		        }
		    }
		}

	    // ActionButton {
	    //     iconSource: "../art/icons/attack.png"
	    //     exclusiveGroup: actionButtons
	    //     onCheckedChanged: {
	    //     	if (checked) {
	    //     		actionBar.parent.right_action = "attack"
	    //     	}
	    //     }
	    // }

	    Item {
	    	height: 32
	    	width: 10
	    }

	    ActionButton {
	        iconSource: "../art/icons/craft.png"
	        onCheckedChanged: {
	        	if (checked) {
	        		crafting_window.visible = true
	        	} else {
	        		crafting_window.visible = false
	        	}
	        }
	        allowNone: true
	    }

	    ActionButton {
	        iconSource: "../art/icons/inventory.png"
	        allowNone: true
	        onCheckedChanged: {
	        	if (checked) {
	        		inventory_window.visible = true
	        	} else {
	        		inventory_window.visible = false
	        	}
	        }
	    }

	    // ActionButton {
	    //     iconSource: "../art/icons/character.png"
	    //     allowNone: true
	    // }

	    ActionButton {
	        iconSource: "../art/icons/society.png"
	        allowNone: true
	        onCheckedChanged: {
	        	if (checked) {
	        		society_window.visible = true
	        	} else {
	        		society_window.visible = false
	        	}
	        }
	    }
	}
}