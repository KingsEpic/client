import QtQuick 2.0

Item {
	id: cspriteObj
	property string simpleName : ""
  property int entityID: 0
  property string name : "None"
  property color onHoverColor: "gold"
  property color borderColor: "white"

  property int tx : 0// Tile x
  property int ty : 0// Tile y

  property int dsx : 0 // Destination x and y in pixel coords
  property int dsy : 0 
  property int dtime : 0 // How long to complete the move

  property bool moving : false
  property int move_type : 0

  state: "NORMAL"

  onSimpleNameChanged: {
    csImage.source = "art/" + simpleName + ".png"
    csImage.y = 32 - csImage.sourceSize.height
    csImage.x = 32 - csImage.sourceSize.width
  }

  Image {
  	id: csImage
    source: "art/blank32.png"
  }

  width: {
    if (csImage.width < 32) {
      return csImage.width
    } else {
      return 32
    }
  }
  height: {
    if (csImage.height < 32) {
      return csImage.height
    } else {
      return 32
    }
  }

  // Rectangle {
  //   width: parent.width
  //      height: parent.height
  //   border.width: 2
  //   border.color: "black"
  // }

	MouseArea {
       id: mouseArea
       acceptedButtons: Qt.LeftButton | Qt.RightButton
       width: parent.width
       height: parent.height
       anchors.margins: -10
       hoverEnabled: true         //this line will enable mouseArea.containsMouse

       onEntered: {
        csImage.scale = 0.8
        cspriteObj.parent.parent.hovered_item = name
      }
       onExited:  csImage.scale = 1.0
       onClicked: {

        if (mouse.button == Qt.LeftButton)
        {
            console.log("Left click, tx, ty, y, z: " + tx + ", " + ty + ", " + cspriteObj.y + ", " + cspriteObj.z + ". simpleName: " + simpleName + ", name: " + name + ". Source: " + csImage.source)
            // game.destroyMe(cspriteObj) // Just for testing!
        }
        else if (mouse.button == Qt.RightButton)
        {
            console.log("Right click, tx, ty: " + tx + ", " + ty)

        	switch (cspriteObj.parent.parent.right_action) {
            case "harvest":
              game.createJob("moveharvest", game.delegated, {"x": tx, "y": ty, "layer": cspriteObj.parent.parent.map_layer})
              break
            case "build":
              game.createJob("movebuild", game.delegated, {"x": tx, "y": ty, "layer": cspriteObj.parent.parent.map_layer, "archetype": game.selectedArchetype})
              break
            case "attacl":
              game.createJob("attack", game.delegated, {"x": tx, "y": ty, "layer": cspriteObj.parent.parent.map_layer, "target": game.selectedEntity})
              break
        		default:
              // game.createJob("walkto", game.delegated, "{\"x\": " + tx + ", \"y\": "+ ty + ", \"layer\": \"" + cspriteObj.parent.parent.map_layer + "\"}")
              game.createJob("walkto", game.delegated, {"x": tx, "y": ty, "layer": cspriteObj.parent.parent.map_layer, "distance": 0})
	        		// game.rightClick(cspriteObj)
	        		break
        	}
        }
      }
	}

	states: [
		State {
			name: "WALKING"
			when: moving == true
		}
	]

	onMovingChanged: {
		if (moving == true) {
			walkAnim.running = true
		} else {
			walkAnim.running = false
		}
	}

	ParallelAnimation {
		id: walkAnim

		running: false

		NumberAnimation { target: cspriteObj; property: "x"; to: dsx; duration: dtime }
		NumberAnimation { target: cspriteObj; property: "y"; to: dsy; duration: dtime }

	}

}