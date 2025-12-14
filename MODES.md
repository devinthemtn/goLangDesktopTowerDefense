# 🎮 Game Modes Guide - Tower Defense

## 🎯 Mode Overview

Tower Defense features two distinct game modes, each offering unique challenges and gameplay experiences. Choose your preferred style of play from the main menu and test your strategic skills!

---

## 🏆 **Mode 1: Normal Mode (Campaign)**
*Strategic progression through designed challenges*

### 🎯 **Core Concept**
Normal Mode is a **structured campaign** featuring 10 carefully designed levels, each with unique objectives, enemy configurations, and strategic challenges. Progress through increasingly difficult scenarios to achieve ultimate victory.

### 📊 **Level Progression System**

#### **Level Structure**
Each level has pre-designed parameters:
- **Fixed enemy count** (5-23 enemies per level)  
- **Scaled enemy health** (50-245 HP)
- **Progressive speed increases** (1.0x - 2.0x speed)
- **Balanced starting resources** ($150-$375)
- **Completion bonuses** ($75-$325)

#### **Difficulty Curve**
```
Level 1: Tutorial        →  5 enemies,  50 HP, 1.0x speed, $150 start
Level 5: Mid-Challenge   → 13 enemies, 115 HP, 1.4x speed, $250 start  
Level 10: Final Battle  → 23 enemies, 245 HP, 1.9x speed, $375 start
```

### 🎮 **Gameplay Mechanics**

#### **Victory Conditions**
- ✅ **Complete all 10 levels** to achieve campaign victory
- ✅ **Eliminate every enemy** in each level to progress
- ✅ **Survive each level** without losing all lives

#### **Progression System**
- 🔄 **Linear advancement** - must complete levels in order
- 💰 **Carry-over resources** - money earned transfers between levels
- 🏗️ **Tower persistence** - placed towers remain for next level
- 📈 **Escalating challenges** - each level introduces new difficulties

#### **Level-Specific Features**
- 📝 **Unique descriptions** for each level's theme and challenge
- ⚡ **Custom spawn timing** optimized for strategic gameplay
- 🎯 **Balanced objectives** requiring different tower strategies
- 🏆 **Completion rewards** encouraging efficient play

### 🎯 **Strategic Considerations**

#### **Resource Management**
- **Early Investment**: Levels 1-3 focus on establishing economy
- **Mid-Game Scaling**: Levels 4-7 require strategic tower upgrades  
- **Late Game Power**: Levels 8-10 demand optimal tower combinations
- **Efficiency Focus**: Limited enemies require cost-effective solutions

#### **Tower Strategy by Phase**
```
Levels 1-3: Basic + Slow Tower foundation
Levels 4-6: Heavy + Sniper Tower upgrades
Levels 7-10: Laser + Splash Tower combinations
```

### 🏆 **Victory & Rewards**
- 🎊 **Campaign Completion**: Unlock victory screen with congratulations
- 📊 **Performance Tracking**: Monitor progress through all 10 levels  
- 🔄 **Replay Value**: Restart campaign to try different strategies
- 🎮 **Perfect Runs**: Challenge yourself to complete without losing lives

---

## ♾️ **Mode 2: Endless Mode**
*Infinite survival against escalating threats*

### 🎯 **Core Concept**
Endless Mode is a **survival challenge** featuring infinite waves of enemies with exponentially increasing difficulty. Test your limits and see how long you can survive against overwhelming odds.

### 📈 **Difficulty Scaling System**

#### **Exponential Growth**
Every wave increases difficulty by **15%**, affecting:
- 🔥 **Enemy Health**: Base 50 HP × (1.15^wave)
- 👥 **Enemy Count**: Base 5 enemies × (1 + wave×0.2)  
- ⚡ **Enemy Speed**: Base 1.0 + (wave×0.05)
- ⏱️ **Spawn Rate**: Faster spawning (minimum 0.3s delay)

#### **Scaling Examples**
```
Wave 1:   5 enemies,  50 HP, 1.0x speed, 2.0s spawn delay
Wave 10: 25 enemies, 203 HP, 1.5x speed, 1.5s spawn delay
Wave 20: 45 enemies, 818 HP, 2.0x speed, 1.0s spawn delay
Wave 50: 105 enemies, 43,542 HP, 3.5x speed, 0.3s spawn delay
```

### 🎮 **Gameplay Mechanics**

#### **Survival Objectives**
- ♾️ **Survive indefinitely** against infinite enemy waves
- 📊 **Track performance** by maximum wave reached
- 🏆 **Personal records** for longest survival runs
- 💪 **Endurance challenge** testing long-term strategy

#### **Wave Progression**
- 🌊 **Continuous waves** with no breaks between levels
- 💰 **Wave bonuses** increase with survival time ($50 + wave×$10)
- 🏗️ **Tower persistence** - all towers remain between waves
- ⚡ **Dynamic adaptation** - adjust strategy as difficulty scales

#### **Scaling Mechanics**
- 📈 **Exponential enemy health** growth (1.15× per wave)
- 👥 **Linear enemy count** increase (+0.2× base per wave)  
- ⚡ **Linear speed** increases (+0.05× per wave)
- ⏱️ **Decreasing spawn delays** (faster enemy arrival)

### 🎯 **Strategic Considerations**

#### **Long-Term Planning**
- **Early Economy**: Invest in cost-effective Basic towers initially
- **Mid-Game Transition**: Switch to high-damage Sniper/Heavy towers
- **Late Game Survival**: Focus on area damage and crowd control
- **Exponential Preparation**: Plan for enemies with massive health pools

#### **Survival Strategy Evolution**
```
Waves 1-10:   Basic tower economy building
Waves 11-25:  Heavy damage tower investment  
Waves 26-50:  Area damage and crowd control focus
Waves 51+:    Optimization and perfect positioning
```

#### **Economic Strategy**
- 💰 **Reinvestment Focus**: Spend wave bonuses on stronger towers
- 🏗️ **Infrastructure Building**: Establish strong defensive lines early
- ⚡ **Efficiency Scaling**: Maximize damage per dollar as costs rise
- 🎯 **Strategic Positioning**: Optimize tower placement for infinite scaling

### 📊 **Performance Tracking**
- 🏆 **Wave Counter**: Track your highest wave reached
- 📈 **Difficulty Display**: Real-time difficulty multiplier shown
- 💰 **Economic Efficiency**: Monitor money management over time
- ⏱️ **Survival Time**: Track total game duration

---

## 🎮 **Mode Selection & Navigation**

### 🖱️ **Main Menu Controls**
- **↑/↓ Arrow Keys**: Navigate between game mode options
- **Enter/Space**: Select highlighted game mode
- **ESC**: Exit game application

### 🎯 **In-Game Controls** (Both Modes)
- **Mouse Click**: Place selected tower at cursor position
- **Keys 1-6**: Select different tower types (Basic, Heavy, Sniper, Laser, Splash, Slow)
- **ESC/P**: Pause current game (access pause menu)
- **M**: Return to main menu (from pause screen)
- **R**: Restart current mode (on game over screen)

### 🎪 **Game State Management**
- 🎯 **Menu State**: Mode selection and navigation
- ▶️ **Playing State**: Active gameplay with full tower controls
- ⏸️ **Paused State**: Game frozen with overlay menu
- 💀 **Game Over State**: Display final stats with restart options
- 🏆 **Victory State**: Campaign completion celebration (Normal Mode only)

---

## 🏆 **Strategic Mode Comparison**

### 📊 **Mode Characteristics**

| Feature | Normal Mode | Endless Mode |
|---------|-------------|--------------|
| **Objective** | Complete 10 levels | Survive infinite waves |
| **Difficulty** | Pre-designed curve | Exponential scaling |
| **Duration** | Fixed campaign length | Unlimited playtime |
| **Strategy** | Level-specific tactics | Long-term adaptation |
| **Victory** | Complete Level 10 | Personal best waves |
| **Resources** | Level-based starting money | Continuous wave bonuses |
| **Challenge** | Designed encounters | Scaling mathematical challenge |
| **Replayability** | Perfect run attempts | High score competition |

### 🎯 **Recommended for Different Players**

#### 🎓 **Normal Mode - Best For:**
- **New players** learning game mechanics
- **Strategic planners** who enjoy designed challenges  
- **Goal-oriented players** who want clear victory conditions
- **Campaign lovers** who prefer structured progression
- **Completionists** who want to "beat the game"

#### 🏃 **Endless Mode - Best For:**
- **Experienced players** seeking maximum challenge
- **Competitive gamers** wanting high score competition
- **Survival enthusiasts** who enjoy escalating difficulty
- **Optimization experts** perfecting long-term strategies
- **Challenge seekers** testing their absolute limits

---

## 🎯 **Advanced Mode Strategies**

### 🏆 **Normal Mode Mastery**
1. **Study level preview** - Use the 3-second level info display
2. **Plan tower progression** - Know which towers work best per level
3. **Economic efficiency** - Spend wisely on each level's enemy count
4. **Perfect runs** - Aim to complete without losing any lives
5. **Speed completion** - Challenge yourself to finish levels quickly

### ♾️ **Endless Mode Mastery**  
1. **Exponential mindset** - Plan for massive enemy health scaling
2. **Economic scaling** - Reinvest wave bonuses efficiently
3. **Area damage focus** - Crowd control becomes essential at high waves
4. **Perfect positioning** - Optimize every tower placement for maximum efficiency
5. **Mathematical optimization** - Calculate optimal tower combinations for scaling

### 🧠 **Meta-Strategy Tips**
- **Mode switching** - Use Normal Mode to practice for Endless runs
- **Tower familiarity** - Master all 6 tower types for optimal performance
- **Pattern recognition** - Learn enemy movement patterns for better placement
- **Resource curves** - Understand when to save vs. spend in each mode
- **Failure analysis** - Learn from defeats to improve future runs

---

## 🎉 **Getting Started**

### 🚀 **First Time Players**
1. **Start with Normal Mode** to learn game mechanics
2. **Complete the first few levels** to understand tower basics  
3. **Experiment with all tower types** to find your preferred strategy
4. **Try Endless Mode** once comfortable with controls
5. **Challenge yourself** to improve personal best scores

### 🎯 **Quick Mode Selection**
```bash
# Launch game
./tower-defense

# At main menu:
# ↓ Select "Normal Mode" → Enter (for campaign)  
# ↓ Select "Endless Mode" → Enter (for survival)
# ↓ Select "Exit Game" → Enter (to quit)
```

Choose your battlefield and prove your strategic mastery! 🏆