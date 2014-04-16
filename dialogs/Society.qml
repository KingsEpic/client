import QtQuick 2.1
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0
// import GoExtensions 1.0


CitWindow {
    id: societyWindow
    width: 350
    height: 350

    Rectangle {
        clip: true
        x: content_x
        y: content_y
        width: content_width
        height: content_height

        Component {
            id: jobDelegate
            Item {
                width: parent.width; height: 40
                Column {
                    Text { text: '<b>Ent ID:</b> ' + jobModel.get(index).workerID }
                    Text { text: '<b>Name:</b> ' + jobModel.get(index).readableName }
                }
            }
        }

        ListView {
            anchors.fill: parent
            model: jobModel.len
            delegate: jobDelegate
            highlight: Rectangle { color: "lightsteelblue"; radius: 5 }
            focus: true
        }
    }
}