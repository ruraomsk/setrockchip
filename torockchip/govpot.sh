export QT_QPA_PLATFORM=eglfs
export QT_QPA_EGLFS_INTEGRATION=eglfs_kms
export QT_QPA_EGLFS_KMS_CONFIG=/etc/qt5/eglfs_kms.json
export QT_QPA_EGLFS_WIDTH=1280
export QT_QPA_EGLFS_HEIGHT=800
export QT_QPA_EGLFS_PHYSICAL_WIDTH=480
export QT_QPA_EGLFS_PHYSICAL_HEIGHT=270
cd /home/rura
while true
do
    ./vpot > /dev/null 2>/dev/null
done 
