import QtQuick 2.1
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0


Item {
    id: rawDialog
    property int margin: 11
    // Rectangle {
    //     width: parent.width
    //     height: parent.height
    //     radius: 10
    //     color: "#cccccc"
    // }

	// Rectangle {
	//     anchors.centerIn: parent
	//     width: parent.width + 6
	//     height: parent.height + 6
	//     radius: 5
	//     opacity: 0.4
	//     color: "#222222"
	// }

	Image {
    	id: actionIcon
    	width: parent.width; height: parent.height;
    	smooth: true
        source: "../art/button.png"
  }
}