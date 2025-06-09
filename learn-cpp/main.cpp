#include <iostream>
#include <string>
#include <string_view>
#include <cmath>

class Ball
{
private:
    std::string m_color{ "black" };
    double m_radius{ 10.0 };
    static constexpr double PI{ 3.14159265358979323846 };

public:
    // Constructors
    Ball(double radius)
        : Ball{ "black", radius } // delegate to the other constructor
    {
        // We don't need to call print() here since it will be called by
        // the constructor we delegate to
    }

    Ball(std::string_view color="black", double radius=10.0)
        : m_color{ color }
        , m_radius{ radius }
    {
        print();
    }

    // Getters
    const std::string& getColor() const { return m_color; }
    double getRadius() const { return m_radius; }

    // Setters
    void setColor(std::string_view color) { m_color = color; }
    void setRadius(double radius) { m_radius = radius; }

    // Mathematical operations
    double getVolume() const
    {
        return (4.0 / 3.0) * PI * std::pow(m_radius, 3);
    }

    double getSurfaceArea() const
    {
        return 4.0 * PI * std::pow(m_radius, 2);
    }

    // Comparison operators
    bool operator==(const Ball& other) const
    {
        return (m_radius == other.m_radius) && (m_color == other.m_color);
    }

    bool operator!=(const Ball& other) const
    {
        return !(*this == other);
    }

    bool operator<(const Ball& other) const
    {
        if (m_radius != other.m_radius)
            return m_radius < other.m_radius;
        return m_color < other.m_color;
    }

    bool operator>(const Ball& other) const
    {
        return other < *this;
    }

    bool operator<=(const Ball& other) const
    {
        return !(other < *this);
    }

    bool operator>=(const Ball& other) const
    {
        return !(*this < other);
    }

    // Binary operators
    Ball operator+(const Ball& other) const
    {
        // Creates a new ball with combined volume
        double combined_volume = this->getVolume() + other.getVolume();
        double new_radius = std::cbrt((3.0 * combined_volume) / (4.0 * PI));
        return Ball(m_color, new_radius);
    }

    // Print function
    void print() const
    {
        std::cout << "Ball(" << m_color << ", " << m_radius << ")\n"
                  << "  Volume: " << getVolume() << "\n"
                  << "  Surface Area: " << getSurfaceArea() << "\n";
    }
};

int main()
{
    // Create balls with different properties
    Ball def{};
    Ball blue{ "blue" };
    Ball twenty{ 20.0 };
    Ball blueTwenty{ "blue", 20.0 };

    // Test comparison
    if (twenty == blueTwenty)
        std::cout << "\nBalls are equal in size and color\n";
    else if (twenty < blueTwenty)
        std::cout << "\nFirst ball is smaller\n";
    else
        std::cout << "\nFirst ball is larger\n";

    // Test combining balls
    Ball combined = twenty + blueTwenty;
    std::cout << "\nCombined ball properties:\n";
    combined.print();

    return 0;
}