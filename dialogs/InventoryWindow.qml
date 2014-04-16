import QtQuick 2.1
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0

CitWindow {
    id: societyWindow
    width: 370
    height: 210

    title: "Inventory"

    Inventory {
        clip: true
        x: content_x
        y: content_y
        width: content_width
        height: content_height
    }

}