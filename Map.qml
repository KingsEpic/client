import QtQuick 2.0
import "dialogs"
// import GoExtensions 1.0

Item {
	id: map
    width: (3000*32) + 16  // Lazy not sending map size yet, so adding 16 so it's obvious if edge of this view is reached
    height: (3000*32) + 16

    focus: true

    onVisibleChanged: {
    	// For some reason I don't need this on Linux, but I do on windows:
    	map.focus = map.visible
    }

    Keys.onPressed: {
    	if (event.isAutoRepeat == false) {
    		if ((event.key == Qt.Key_A) || (event.key == Qt.Key_Left)) {
    			leftAnim.start()
    			rightAnim.stop()
    		}
    		if ((event.key == Qt.Key_D) || (event.key == Qt.Key_Right)) {
    			rightAnim.start()
    			leftAnim.stop()
    		}
    		if ((event.key == Qt.Key_S) || (event.key == Qt.Key_Down)) {
    			downAnim.start()
    			upAnim.stop()
    		}
    		if ((event.key == Qt.Key_W) || (event.key == Qt.Key_Up)) {
    			upAnim.start()
    			downAnim.stop()
    		}
	    	// console.log("Key " + event.key + " pressed.  Count: " + event.count + ". Modifiers: " + event.modifiers)
	    }
    }

    Keys.onReleased: {
    	if (event.isAutoRepeat == false) {
    		if ((event.key == Qt.Key_A) || (event.key == Qt.Key_Left)) {
    			leftAnim.stop()
    			x = Math.round(x); // Makes sure the map stops on whole numbers so pixels aren't blurred
    		}
    		if ((event.key == Qt.Key_D) || (event.key == Qt.Key_Right)) {
    			rightAnim.stop()
    			x = Math.round(x);
    		}
    		if ((event.key == Qt.Key_S) || (event.key == Qt.Key_Down)) {
    			downAnim.stop()
    			y = Math.round(y);
    		}
    		if ((event.key == Qt.Key_W) || (event.key == Qt.Key_Up)) {
    			upAnim.stop()
    			y = Math.round(y);
    		}
		    // console.log("Key " + event.key + "released.  Count: " + event.count)
		}
    }

    NumberAnimation on x {
        id: leftAnim
        running: false
        to: map.x + 2000
        duration: 5000
    }

    NumberAnimation on x {
        id: rightAnim
        running: false
        to: map.x - 2000
        duration: 5000
    }
    NumberAnimation on y {
        id: upAnim
        running: false
        to: map.y + 2000
        duration: 5000
    }

    NumberAnimation on y {
        id: downAnim
        running: false
        to: map.y - 2000
        duration: 5000
    }

}