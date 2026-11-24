#pragma once
#include "include/configs/common/Outbound.h"

namespace Configs {
    // The built-in Tor exit; the core pins the executable and data paths itself.
    class tor : public outbound {
    public:
        static QString ExecutablePath();
        static bool IsInstalled();

        QString DisplayType() override;

        bool HasServerEndpoint() override { return false; }

        BuildResult Build() override;
    };
}
