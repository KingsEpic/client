import QtQuick 2.0
import "dialogs"

View {
	id: mapView
	width: parent.width; height: parent.height;
	x: 0
	y: 0

    property string map_layer: "terra"

	property int ui_z: 1000000
    property string right_action
    property string hovered_item


    Component.onCompleted: {
        game.map = map
        game.windowZ = ui_z
    }

    ActionBar {
        z: ui_z

        width: parent.width; height: parent.height;
    }

    Society {
        id: society_window
        x: width / 2
        y: width / 2
        z: ui_z
        title: "Society"
        visible: false
    }

    InventoryWindow {
        id: inventory_window
        x: width / 2
        y: width / 2
        z: ui_z
        visible: false
    }

    // Hotbar {
    //     id: hotbar_window
    //     height: 32
    //     width: 32*10
    //     anchors.bottom: parent.bottom
    //     anchors.horizontalCenter: parent.horizontalCenter
    //     z: ui_z
    // }

    Crafting {
        id: crafting_window
        x: 64
        y: 64
        z: ui_z
        visible: false
    }

    Rectangle {
    	id: statusRect
    	x: 8
    	y: statusRect.parent.height - 24
    	z: ui_z

        width: 128
        height: 16
        radius: 10
        color: "#cccccc"
	    Text {
	    	y: 0
	    	x: 8
	    	color: "black"
	    	text: hovered_item
	    }
    }

	Map {
		id: map
	}

}